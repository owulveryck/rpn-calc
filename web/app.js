(function () {
    "use strict";

    var inputBuffer = "";
    var wasmReady = false;

    async function loadWasm() {
        var go = new Go();
        var result = await WebAssembly.instantiateStreaming(
            fetch("rpn.wasm"),
            go.importObject
        );
        go.run(result.instance);
        wasmReady = true;

        var saved = localStorage.getItem("rpn_state");
        if (saved) {
            rpnRestore(saved);
        }
        render(rpnGetState());
    }

    function render(stateJSON) {
        var state = JSON.parse(stateJSON);
        var lines = document.getElementById("stack-lines");
        lines.innerHTML = "";

        var stack = state.stack || [];
        var start = Math.max(0, stack.length - 8);
        for (var i = start; i < stack.length; i++) {
            var level = stack.length - i;
            var entry = document.createElement("div");
            entry.className = "stack-entry";
            var lbl = document.createElement("span");
            lbl.className = "stack-level";
            lbl.textContent = level + ":";
            var val = document.createElement("span");
            val.className = "stack-value";
            val.textContent = stack[i];
            entry.appendChild(lbl);
            entry.appendChild(val);
            lines.appendChild(entry);
        }

        if (stack.length === 0) {
            var empty = document.createElement("div");
            empty.className = "stack-entry";
            var emptyLbl = document.createElement("span");
            emptyLbl.className = "stack-level";
            emptyLbl.textContent = "1:";
            var emptyVal = document.createElement("span");
            emptyVal.className = "stack-value";
            emptyVal.textContent = "";
            empty.appendChild(emptyLbl);
            empty.appendChild(emptyVal);
            lines.appendChild(empty);
        }

        renderHistory(state.history || []);

        document.getElementById("angle-mode").textContent = state.angleMode || "DEG";
        document.getElementById("base-mode").textContent = state.baseMode || "DEC";
        document.getElementById("input-buffer").textContent = inputBuffer;

        saveState();
    }

    function renderHistory(history) {
        var container = document.getElementById("history-lines");
        container.innerHTML = "";
        for (var i = 0; i < history.length; i++) {
            var entry = document.createElement("div");
            entry.className = "history-entry";
            entry.textContent = history[i].expr;
            container.appendChild(entry);
        }
        var historyPanel = document.getElementById("history-display");
        historyPanel.scrollTop = historyPanel.scrollHeight;
    }

    var saveTimeout = null;
    function saveState() {
        if (saveTimeout) clearTimeout(saveTimeout);
        saveTimeout = setTimeout(function () {
            if (wasmReady) {
                localStorage.setItem("rpn_state", rpnSerialize());
            }
        }, 500);
    }

    function pushInput() {
        if (inputBuffer === "") return false;
        var result = rpnEnter(inputBuffer);
        inputBuffer = "";
        render(result);
        return true;
    }

    function handleDigit(d) {
        if (!wasmReady) return;
        if (d === "." && inputBuffer.indexOf(".") !== -1) return;
        inputBuffer += d;
        document.getElementById("input-buffer").textContent = inputBuffer;
    }

    function handleAction(action) {
        if (!wasmReady) return;

        if (action === "ENTER") {
            if (inputBuffer !== "") {
                pushInput();
            } else {
                render(rpnExecute("DUP"));
            }
            return;
        }

        if (action === "BACKSPACE") {
            if (inputBuffer.length > 0) {
                inputBuffer = inputBuffer.slice(0, -1);
                document.getElementById("input-buffer").textContent = inputBuffer;
            } else {
                render(rpnExecute("DROP"));
            }
            return;
        }

        if (action === "UNDO") {
            render(rpnUndo());
            return;
        }

        if (action === "LAST") {
            render(rpnLast());
            return;
        }

        if (action === "NEG") {
            if (inputBuffer.length > 0) {
                if (inputBuffer.charAt(0) === "-") {
                    inputBuffer = inputBuffer.slice(1);
                } else {
                    inputBuffer = "-" + inputBuffer;
                }
                document.getElementById("input-buffer").textContent = inputBuffer;
                return;
            }
            render(rpnExecute("NEG"));
            return;
        }

        if (action === "EEX") {
            if (inputBuffer !== "" && inputBuffer.indexOf("e") === -1) {
                inputBuffer += "e";
                document.getElementById("input-buffer").textContent = inputBuffer;
            }
            return;
        }

        // For operators, push pending input first
        if (["+", "-", "*", "/", "Y^X", "SWAP", "MIN", "MAX"].indexOf(action) !== -1) {
            pushInput();
        }

        render(rpnExecute(action));
    }

    // Button click handlers
    document.getElementById("button-grid").addEventListener("click", function (e) {
        var btn = e.target.closest("button");
        if (!btn) return;

        var digit = btn.getAttribute("data-digit");
        if (digit !== null) {
            handleDigit(digit);
            return;
        }

        var action = btn.getAttribute("data-action");
        if (action !== null) {
            handleAction(action);
        }
    });

    // Keyboard handler
    document.addEventListener("keydown", function (e) {
        if (!wasmReady) return;

        var key = e.key;
        var btn = null;

        if (key >= "0" && key <= "9") {
            handleDigit(key);
            btn = document.querySelector('[data-digit="' + key + '"]');
        } else if (key === ".") {
            handleDigit(".");
            btn = document.querySelector('[data-digit="."]');
        } else if (key === "Enter") {
            e.preventDefault();
            handleAction("ENTER");
            btn = document.getElementById("btn-enter");
        } else if (key === "+" || key === "-" || key === "*" || key === "/") {
            e.preventDefault();
            handleAction(key);
            btn = document.querySelector('[data-action="' + key + '"]');
        } else if (key === "Backspace") {
            e.preventDefault();
            handleAction("BACKSPACE");
        } else if (key === "Delete") {
            e.preventDefault();
            handleAction("CLEAR");
        } else if (key === "Escape") {
            inputBuffer = "";
            document.getElementById("input-buffer").textContent = "";
        } else if (key === "s" || key === "S") {
            handleAction("SWAP");
        } else if (key === "d" || key === "D") {
            handleAction("DUP");
        } else if (key === "u" || key === "U" || (e.ctrlKey && key === "z")) {
            e.preventDefault();
            handleAction("UNDO");
        } else if (key === "^") {
            pushInput();
            handleAction("Y^X");
        }

        if (btn) {
            btn.classList.add("key-active");
            setTimeout(function () {
                btn.classList.remove("key-active");
            }, 100);
        }
    });

    // Prevent zoom on double-tap
    document.addEventListener("touchend", function (e) {
        if (e.target.closest("button")) {
            e.preventDefault();
        }
    }, { passive: false });

    if ("serviceWorker" in navigator) {
        navigator.serviceWorker.register("sw.js");
    }

    loadWasm();
})();
