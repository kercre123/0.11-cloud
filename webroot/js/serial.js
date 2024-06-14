/*
type ESNValidRequest struct {
	ESN string `json:"esn"`
}

type ESNValidResponse struct {
	IsValid bool `json:"esn_isvalid"`
	IsNew   bool `json:"esn_isnew"`
}
*/

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
            }
            
        }
    )
}

function goToSettingsPage(esn) {
    window.location.href = "./settings.html?serial=" + esn;
}