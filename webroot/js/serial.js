/*
type ESNValidRequest struct {
	ESN string `json:"esn"`
}

type ESNValidResponse struct {
	IsValid   bool `json:"esn_isvalid"`
	IsNew     bool `json:"esn_isnew"`
	MatchesIP bool `json:"matches_ip"`
}
*/

document.getElementById("serialInput").addEventListener("keypress", function(event) {
    if (event.key === "Enter") {
      event.preventDefault();
      document.getElementById("doSerialInput").click();
    }
});

function sendSerialInput() {
    var esnInput = document.getElementById("serialInput").value
    fetch("/api/is_esn_valid", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: '{"esn":"' + esnInput + '"}'
    })
        .then((response) => response.text())
        .then((response) => {
            var resp = JSON.parse(response)
            if (!resp["esn_isvalid"]) {
                alert("This ESN is not valid.")
                return
            } else if (!resp["matches_ip"]) {
                alert("This ESN does not match your IP address. It is possible your IP has changed. Try a voice command then try this again.")
                return
            } else {
                goToSettingsPage(esnInput)
            }
        }
    )
}

function goToSettingsPage(esn) {
    window.location.href = "./settings.html?serial=" + esn;
}