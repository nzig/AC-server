function sendData(){
    let temp = document.getElementById('temp').value;
    let formData = new FormData();
    formData.append('temp', temp);
    fetch('/send', {
        method: 'POST',
        credentials: 'include',
        body: formData
    }).then(r => r.text()).then(t => {
        let out =  document.getElementById('output');
        out.innerText = t;
        out.classList.add('in');
    })
}