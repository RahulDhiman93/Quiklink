document.getElementById("shortenButton").addEventListener("click", function () {
    const originalUrl = document.getElementById("originalUrl").value
    const alias = document.getElementById("alias").value
    const copyButton = document.getElementById("copyButton");

    if (!isURL(originalUrl)) {
        copyButton.style.display = "none"
        document.getElementById("originalUrl").value = "";
        document.getElementById("shortenedUrl").href = "#!";
        document.getElementById("shortenedUrl").textContent = "";
        alert("Enter a valid URL");
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
                alert(data["message"]);
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
        })
        .catch(error => {
            alert(error);
            window.location.reload();
        });
});

document.getElementById("openUrlButton").addEventListener("click", function () {
    console.log("Event listener triggered");
    console.log("INSIDE OPEN URL")
    const shortenedUrl = document.getElementById("shortenedUrlBox");
    console.log("Shortened URL:", shortenedUrl.value); // Add this line
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
    alert("Copied to clipboard");
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

