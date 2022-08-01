function upload(button) {
    const messagesDiv = document.getElementById("messages")
    messagesDiv.innerHTML = ""
    button.disabled = true

    formData = new FormData();           
    formData.append("CHINESE_FILE", CHINESE_FILE.files[0]);
    formData.append("ORIGINAL_FILE", ORIGINAL_FILE.files[0]);
    formData.append("REVISED_FILE", REVISED_FILE.files[0]);

    fetch('/comparefiles', {
        method: "POST", 
        body: formData
    })
    .then(response => {
        if (response.status != 200) {
            throw new Error("Sorry, there was a problem trying to compare the files. Make sure you have uploaded files correctly.")
        }
        return response.blob()
    })
    .then(blob => {
        download(blob, getNameForFile("output", "xlsx"))
    })
    .catch(err => messagesDiv.innerHTML = err)
    .finally(() => button.disabled = false);
}

function getNameForFile(prefix, extension) {
    d = new Date()
    filename = [
        d.getFullYear(),
        pad(d.getMonth()+1, 2, '0'),
        pad(d.getDate(), 2, '0'),
        pad(d.getHours(), 2, '0'),
        pad(d.getMinutes(), 2, '0'),
        pad(d.getSeconds(), 2, '0')
    ].join('-')
    return prefix + "_" + filename + "." + extension
}

function download(blob, filename) {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = url;
    a.download = filename;

    document.body.appendChild(a);
    a.click();
    
    a.remove();
    window.URL.revokeObjectURL(url);
  }

  function pad(num, padlen, padchar) {
    var pad_char = typeof padchar !== 'undefined' ? padchar : '0';
    var pad = new Array(1 + padlen).join(pad_char);
    return (pad + num).slice(-pad.length);
}