"use strict";

async function sendData(){
    statusElement.show(false);

    try {
        const temperature = document.getElementById('temp').value;
        const formData = new FormData();
        formData.append('temp', temp);
        const r = await fetch('/send', {
            method: 'POST',
            credentials: 'include',
            body: formData
        });
        if(!r.ok)
        {
            throw new Error(`Posting to server failed: ${r.statusText}`);
        }
        const text = await r.text();
        statusElement
            .type('success')
            .message(text)
            .show();
    } catch(e) {
        statusElement
            .type('danger')
            .message(`Exception: ${e.toString()}`)
            .show();
    }
}

let statusElement = null;

class StatusElement {
    constructor() {
        this.el = document.getElementById('output');
        this.currentType = null;
    }

    type(type) {
        if (this.currentType !== null) {
            this.el.classList.remove('alert-' + this.currentType)
        }
        this.el.classList.add('alert-' + type);
        return this;
    }

    message(message) {
        this.el.innerHTML = message;
        return this;
    }

    show(show=true) {
        const cssForShow = 'in';
        if(show) {
            this.el.classList.add(cssForShow);
        } else {
            this.el.classList.remove(cssForShow);
        }
        return this;
    }
}

document.addEventListener("DOMContentLoaded", () => {
    statusElement = new StatusElement();
    document.getElementById("button").addEventListener("click", sendData);
})