function checkAuthorized() {
    fetch("/api/check_if_whitelisted", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then((response) => response.text())
        .then((response) => {
            var resp = JSON.parse(response)
            if (!resp["whitelisted"]) {
                alert("Your IP is still not whitelisted. If you did a voice command, make sure you are on the same network as Vector and that you aren't using mobile data.")
            } else {
                window.location.href = "/"
            }
        }
    )
}