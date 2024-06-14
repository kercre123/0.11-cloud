var urlParams = new URLSearchParams(window.location.search);
var esn = urlParams.get("serial");
getSettings();

/*
        <label for="locationInput">Location (for weather command)</label>
        <input id="locationInput" type="text">
        <br />
        <label for="kgServiceInput">Knowledge Graph Service</label>
        <select name="kgServiceInput" id="kgServiceInput" onchange="checkKG()">
        <option value="" selected>None</option>
        <option value="openai">OpenAI</option>
        <option value="houndify">Houndify</option>
        <option value="together">Together</option></select>
        <br />
        <span id="kgKey" style="display: none">
            <label for="kgKeyInput">KG API Key</label>
            <input id="kgKeyInput" type="text">
        </span>
        <span id="kgClientKey" style="display: none">
            <label for="kgClientKeyInput">Houndify Client Key</label>
            <input id="kgClientKeyInput" type="text">
        </span>
        <span id="kgPrompt" style="display: none">
            <label for="kgPromptInput">LLM Prompt (leave blank for default)</label>
            <input id="kgPromptInput" type="text">
        </span>
*/

function hideE(element) {
    document.getElementById(element).style.display = "none"
}

function showE(element) {
    document.getElementById(element).style.display = "block"
}

function setE(element, set) {
    document.getElementById(element).value = set
}

function checkKG() {
    kgService = document.getElementById("kgServiceInput").value
    switch(kgService) {
        case "openai":
            showE("kgKey")
            hideE("kgClientKey")
            showE("kgPrompt")
            break
        case "together":
            showE("kgKey")
            hideE("kgClientKey")
            showE("kgPrompt")
            break
        case "houndify":
            showE("kgKey")
            showE("kgClientKey")
            hideE("kgPrompt")
            break
        default:
            hideE("kgKey")
            hideE("kgClientKey")
            hideE("kgPrompt")
      }
}

function setValues(location, kgEnabled, kgService, kgKey, kgClientKey, kgPrompt) {
    setE("locationInput", location)
    if (kgEnabled) {
        setE("kgServiceInput", kgService)
        setE("kgKeyInput", kgKey)
        setE("kgClientKeyInput", kgClientKey)
        setE("kgPromptInput", kgPrompt)
    } else {
        setE("kgServiceInput", "")
    }
}

function getE(element) {
    var value = document.getElementById(element).value
    if (!value) {
        return ""
    } else {
        return value
    }
}

function setSettings() {
    var isKGEnabled = false
    if (getE("kgServiceInput") != "none") {
        isKGEnabled = true
    }
    let kgData = 
    {
        enabled: isKGEnabled,
        prompt: getE("kgPromptInput"),
        service: getE("kgServiceInput"),
        apikey: getE("kgKeyInput"),
        clientkey: getE("kgClientKeyInput")
    };
    let settingsData = 
    { 
        esn: esn, 
        location: getE("locationInput"),
        kg: kgData
    };
    let settingsJSON = JSON.stringify(settingsData)
    console.log(settingsJSON)
    fetch("/api/set_settings", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: settingsJSON
    })
        .then((response) => {
            if (!response.ok) {
                alert("There was an error with the settings.")
                console.log(response.text())
            } else {
                getSettings()
                alert("Settings successfully set.")
            }
        })
}

function getSettings() {
    fetch("/api/get_settings", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: '{"esn":"' + esn + '"}'
    })
        .then((response) => response.text())
        .then((response) => {
            var resp = JSON.parse(response)
            if (resp["kg"]["enabled"]) {
                setValues(resp["location"], true, resp["kg"]["service"], resp["kg"]["apikey"], resp["kg"]["clientkey"], resp["kg"]["prompt"])
                checkKG()
            } else {
                setValues(resp["location"], false, "", "", "", "")
                checkKG()
            }
        }
    )
}