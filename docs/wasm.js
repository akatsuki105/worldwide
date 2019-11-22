window.onload = async () => {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch("./gbc.wasm"), go.importObject);
    go.run(result.instance);

    let canvas = document.querySelector("canvas");
    let ctx = canvas.getContext("2d");
    let image = ctx.createImageData(160, 144);
    ctx.scale(1.40625, 1.421);

    let canvasInvisible = document.createElement('canvas');
    canvasInvisible.width = 160;
    canvasInvisible.height = 144;
    let ctx2 = canvasInvisible.getContext('2d');

    const init = async (event) => {
        let files;
        let reader = new FileReader();

        if (event.target.files) {
            files = event.target.files;
        } else {
            files = event.dataTransfer.files;
        }

        reader.onload = (event) => {
            const romData = new Uint8Array(reader.result);
            console.log(romData);
            let gb = new GB(new Uint8Array(romData));

            const frame = () => {
                if (!gb) return;

                let arr = new Uint8Array(160 * 144 * 4);
                gb.next(arr);
                image.data.set(arr);

                ctx2.putImageData(image, 0, 0);
                ctx.drawImage(canvasInvisible, 0, 0);
                window.requestAnimationFrame(frame);
            };

            frame();

            const onKeyDown = (e) => {
                switch (e.key) {
                    case "z":
                        return gb.keyDown("B");
                    case "x":
                        return gb.keyDown("A");
                    case "Shift":
                        return gb.keyDown("Select");
                    case "Enter":
                        return gb.keyDown("Start");
                    case "ArrowLeft":
                        return gb.keyDown("Left");
                    case "ArrowUp":
                        return gb.keyDown("Up");
                    case "ArrowRight":
                        return gb.keyDown("Right");
                    case "ArrowDown":
                        return gb.keyDown("Down");
                }
            };

            const onKeyUp = (e) => {
                switch (e.key) {
                    case "z":
                        return gb.keyUp("B");
                    case "x":
                        return gb.keyUp("A");
                    case "Shift":
                        return gb.keyUp("Select");
                    case "Enter":
                        return gb.keyUp("Start");
                    case "ArrowLeft":
                        return gb.keyUp("Left");
                    case "ArrowUp":
                        return gb.keyUp("Up");
                    case "ArrowRight":
                        return gb.keyUp("Right");
                    case "ArrowDown":
                        return gb.keyUp("Down");
                }
            };

            window.addEventListener("keydown", onKeyDown);
            window.addEventListener("keyup", onKeyUp);
        }

        if (files[0]) {
            reader.readAsArrayBuffer(files[0]);
        }
    }

    const trial = async (event) => {
        const rom = await fetch("./tobu.gb");
        let buf = await rom.arrayBuffer();
        let gb = new GB(new Uint8Array(buf));

        const frame = () => {
            if (!gb) return;

            let arr = new Uint8Array(160 * 144 * 4);
            gb.next(arr);
            image.data.set(arr);

            ctx2.putImageData(image, 0, 0);
            ctx.drawImage(canvasInvisible, 0, 0);
            window.requestAnimationFrame(frame);
        };

        frame();

        const onKeyDown = (e) => {
            switch (e.key) {
                case "z":
                    return gb.keyDown("B");
                case "x":
                    return gb.keyDown("A");
                case "Shift":
                    return gb.keyDown("Select");
                case "Enter":
                    return gb.keyDown("Start");
                case "ArrowLeft":
                    return gb.keyDown("Left");
                case "ArrowUp":
                    return gb.keyDown("Up");
                case "ArrowRight":
                    return gb.keyDown("Right");
                case "ArrowDown":
                    return gb.keyDown("Down");
            }
        };

        const onKeyUp = (e) => {
            switch (e.key) {
                case "z":
                    return gb.keyUp("B");
                case "x":
                    return gb.keyUp("A");
                case "Shift":
                    return gb.keyUp("Select");
                case "Enter":
                    return gb.keyUp("Start");
                case "ArrowLeft":
                    return gb.keyUp("Left");
                case "ArrowUp":
                    return gb.keyUp("Up");
                case "ArrowRight":
                    return gb.keyUp("Right");
                case "ArrowDown":
                    return gb.keyUp("Down");
            }
        };

        window.addEventListener("keydown", onKeyDown);
        window.addEventListener("keyup", onKeyUp);
    }

    // let inputROM = document.getElementById("inputROM");
    // inputROM.addEventListener("change", event => {
    //     inputROM.blur();
    //     return init(event);
    // }, false);

    // let trialROM = document.getElementById("trialROM");
    // trialROM.addEventListener("click", event => {
    //     trialROM.blur();
    //     return trial(event);
    // }, false);

    trial(event);
}