document.getElementById("shortenButton").addEventListener("click", function () {
    const originalUrl = document.getElementById("originalUrl").value
    const alias = document.getElementById("alias").value
    const copyButton = document.getElementById("copyButton");

    if (!isURL(originalUrl)) {
        prompt().error({
            title: "OOPS!!!",
            msg: "Enter a valid URL"
        });
        copyButton.style.display = "none"
        document.getElementById("originalUrl").value = "";
        document.getElementById("shortenedUrl").href = "#!";
        document.getElementById("shortenedUrl").textContent = "";
        return
    }

    const url = '/shorten';
    const data = {
        long_url: originalUrl,
        url_alias: alias
    };

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
        .then(response => response.json())
        .then(data => {
            if (data["ok"] === false) {
                prompt().error({
                    title: "OOPS!!!",
                    msg: data["message"]
                });
                window.location.reload();
                return
            }

            const shortenedUrl = data["short_url"];
            copyButton.style.display = "inline"
            document.getElementById("longUrlBox").value = originalUrl;
            document.getElementById("shortenedUrlBox").value = shortenedUrl;

            //setQrImage(data["qrcode"])
            const urlContainer = document.querySelector('.url-container');
            urlContainer.style.display = "none";

            const shortenedContainer = document.querySelector('.shortened-container');
            shortenedContainer.classList.add('show');
            shortenedContainer.style.display = "block";

            setQrImage(data["qrcode"])
        })
        .catch(error => {
            prompt().error({
                title: "OOPS!!!",
                msg: error
            });
            window.location.reload();
        });
});

document.getElementById("openUrlButton").addEventListener("click", function () {
    const shortenedUrl = document.getElementById("shortenedUrlBox");
    window.open(shortenedUrl.value, '_blank');
});

document.getElementById("copyButton").addEventListener("click", function () {
    const shortenedUrl = document.getElementById("shortenedUrlBox");
    const textArea = document.createElement("textarea");
    textArea.value = shortenedUrl.value;
    document.body.appendChild(textArea);
    textArea.select();
    document.execCommand("copy");
    document.body.removeChild(textArea);
    prompt().success({
        title: "Share your Quiklink",
        msg: "Copied to clipboard"
    });
});

document.getElementById("shortenAgainButton").addEventListener("click", function () {
    document.getElementById("originalUrl").value = ""

    const urlContainer = document.querySelector('.url-container');
    urlContainer.style.display = "block";

    const shortenedContainer = document.querySelector('.shortened-container');
    shortenedContainer.style.display = "none";
});

function isURL(text) {
    if (!text.includes("http") || !text.includes("https")) {
        return false
    }
    const urlPattern = /^(https?:\/\/)?([\w.-]+)\.([a-zA-Z]{2,})(\S*)$/;
    return urlPattern.test(text);
}

function setQrImage(data) {
    const bytes = data;
    const canvasContainer = document.getElementById("canvas-container");
    const canvas = document.getElementById("qrImageCanvas");
    const ctx = canvas.getContext("2d");
    const image = new Image();
    image.src = "data:image/png;base64," + bytes;
    image.onload = function() {
        ctx.drawImage(image, 0, 0);
        canvasContainer.style.display = "block";
    };
}

function prompt() {
    let myToast = function (c) {

        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {

        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
            confirmButtonColor: "#800160",
        })
    }

    let error = function (c) {

        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
            confirmButtonColor: "#800160",
        })
    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c

        const {value: result} = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen() !== undefined) {
                    c.willOpen()
                }
            },
            didOpen: () => {
                if (c.didOpen() !== undefined) {
                    c.didOpen()
                }
            },
        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result)
                    }
                } else {
                    c.callback(false)
                }
            } else {
                c.callback(false)
            }
        }

    }

    return {
        myToast: myToast,
        success: success,
        error: error,
        custom: custom,
    }
}
