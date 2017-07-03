function sendData(){
    let temp = document.getElementById('temp').value;
    fetch('/send', {
        method: 'POST',
        credentials: 'include',
        body: temp
    }).then(r => r.text()).then(t => {
        document.getElementById('output').innerText = t;
    })
}