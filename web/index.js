/// mmmmm drinking a nice spanish tempranillio while writing this shit

const result = document.getElementById('result');
const guessIndicator = document.getElementById('guessIndicator');
const mapImage = document.getElementById('mapImage');
const guessText = document.getElementById('guessText');
const guessBtn = document.getElementById('guessBtn');
const controlRow = document.getElementById('controlRow');

const guesses = 6;

// I yoinked disssss offffa da internettttttt
function uuidv4() {
    return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, c =>
        (+c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> +c / 4).toString(16)
    );
}

let getSession = (forceRefresh = false) => {
    let lastSessionGuid = (forceRefresh ? null : localStorage.getItem('session')) ?? function wow() {
        newSessionId = uuidv4();
        localStorage.setItem('session', newSessionId);
        return newSessionId;
    }();
    return lastSessionGuid;
};

let lastSessionGuid = getSession();

// also got offa da internetz
function humanReadableTimeDiff(diff) {
    let ss = Math.floor(diff / 1000) % 60;
    ss = ss.toString().padStart(2, "0")
    let mm = Math.floor(diff / 1000 / 60) % 60;
    mm = mm.toString().padStart(2, "0")
    let hh = Math.floor(diff / 1000 / 60 / 60);
    hh = hh.toString().padStart(2, "0")

    return `${hh}:${mm}:${ss}`
}


// idk it's pure javascript leave me alone
let getCurrentChallengeImage = () => {
    let x = new XMLHttpRequest();
    x.open('GET', 'https://ss14mapdle-api.croil.net/challenge/' + lastSessionGuid);
    // x.responseType = 'blob';
    x.addEventListener('readystatechange', function () {
        if (this.readyState == 4) {
            switch (this.status) {
                case 200:
                    sessionInfo = JSON.parse(this.response);
                    mapImage.setAttribute('src', 'https://ss14mapdle-api.croil.net/challenge/map/' + lastSessionGuid + '?t=' + Date.now());
                    guessIndicator.innerHTML = '';
                    for (i = 0; i < guesses - 1; i++) {
                        let elem = document.createElement('span');

                        if (i < sessionInfo.session.zoom_level - 1) {
                            elem.innerText = '🟥';
                        } else if (i == sessionInfo.session.zoom_level - 1 && sessionInfo.session.correct === true) {
                            elem.innerText = '🟩';
                        } else {
                            elem.innerText = '⬛';
                        }

                        guessIndicator.appendChild(elem);
                    }
                    if (sessionInfo.session.zoom_level >= guesses || sessionInfo.session.correct === true) {
                        if (sessionInfo.session.correct === true) {
                            result.setAttribute('class', 'card-panel green center-align')
                            result.innerText = sessionInfo.message;
                            result.style['visibility'] = 'visible';
                        }
                        window.setInterval(() => {
                            let diff = Date.parse(sessionInfo.expires_at) - Date.now();
                            controlRow.innerText = `Next challenge in ${humanReadableTimeDiff(diff)}`;
                        }, 1000);
                        let diff = Date.parse(sessionInfo.expires_at) - Date.now();
                        controlRow.innerText = `Next challenge in ${humanReadableTimeDiff(diff)}`;
                    }
                    break;
                case 410: // gone, new challenge = get new challenge thing whatever
                    lastSessionGuid = getSession(true);
                    getCurrentChallengeImage();
            }
        }
    });
    x.send();
};

getCurrentChallengeImage();

const options = {
    "bagel": null,
    "box": null,
    "cache": null,
    "elkridge": null,
    "exo": null,
    "fland": null,
    "marathon": null,
    "oasis": null,
    "oasis brig": null,
    "packed": null,
    "plasma": null,
    "reach": null,
    "saltern": null,
    "serpentcrest": null,
    "snowball": null,
};

document.addEventListener('DOMContentLoaded', function () {
    var elems = document.querySelectorAll('.autocomplete');
    var instances = M.Autocomplete.init(elems, {
        data: options
    });
});

let validateGuess = () => {
    const guess = guessText.value.toLowerCase();

    if (!(guess in options)) {
        result.setAttribute('class', 'card-panel green center-align')
        result.innerText = `${guess} is not a valid map`
        result.style['visibility'] = 'visible';
        return;
    }


    let x = new XMLHttpRequest();
    x.open('POST', 'https://ss14mapdle-api.croil.net/guess/' + lastSessionGuid)
    x.setRequestHeader('content-type', 'application/json')

    x.addEventListener('readystatechange', function () {
        if (this.readyState == 4) {
            guessText.value = '';
            switch (this.status) {
                case 200:
                    var data = JSON.parse(this.response)
                    if (data.correct === true) {
                        result.setAttribute('class', 'card-panel green center-align')
                    } else {
                        result.setAttribute('class', 'card-panel blue center-align')
                    }
                    result.innerText = data.message;
                    result.style['visibility'] = 'visible';
                    getCurrentChallengeImage();
                    break;
            }
        }
    });

    x.send(JSON.stringify({
        "name": guess
    }));

};

guessBtn.addEventListener('click', validateGuess);

guessText.addEventListener('keypress', (ev) => {
    if (ev.keyCode === 13) {
        validateGuess();
    }
});

