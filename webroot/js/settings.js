var urlParams = new URLSearchParams(window.location.search);
var esn = urlParams.get("serial");
getSettings();

var HoundClientID = ""
var HoundClientKey = ""
var OpenAIKey = ""
var TogetherKey = ""
var AiPrompt = ""

var inputs = ["locationInput", "kgKeyInput", "kgClientKey", "kgPromptInput"]

for (var element of inputs) {
    document.getElementById(element).addEventListener("keypress", function(event) {
        if (event.key === "Enter") {
            event.preventDefault();
            document.getElementById("doSetSettings").click();
        }
    });
}

function hideE(element) {
    document.getElementById(element).style.display = "none"
}

function showE(element) {
    document.getElementById(element).style.display = "block"
}

function setE(element, set) {
    document.getElementById(element).value = set
}

function setIHTML(element, set) {
    document.getElementById(element).innerHTML = set
}

function checkKG() {
    kgService = document.getElementById("kgServiceInput").value
    switch(kgService) {
        case "openai":
            setIHTML("kgKeyLabel", "OpenAI API Key")
            setE("kgKeyInput", OpenAIKey)
            setE("kgPrompt", AiPrompt)
            showE("kgKey")
            hideE("kgClientKey")
            showE("kgPrompt")
            break
        case "together":
            setIHTML("kgKeyLabel", "Together API Key")
            setE("kgKeyInput", TogetherKey)
            setE("kgPrompt", AiPrompt)
            showE("kgKey")
            hideE("kgClientKey")
            showE("kgPrompt")
            break
        case "houndify":
            setIHTML("kgKeyLabel", "Houndify Client ID")
            setE("kgKeyInput", HoundClientID)
            setE("kgClientKeyInput", HoundClientKey)
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
    /*
    var HoundClientID
    var HoundClientKey
    var OpenAIKey
    var TogetherKey
    var AiPrompt
    */
    setE("locationInput", location)
    if (kgEnabled) {
        if (kgService == "houndify") {
            HoundClientID = kgKey
            HoundClientKey = kgClientKey
        } else {
            if (kgService == "openai") {
                OpenAIKey = kgKey
            } else if (kgService == "together") {
                TogetherKey = kgKey
            }
            AiPrompt = kgPrompt
        }
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