function upload() {
    formData = new FormData();           
    formData.append("CHINESE_FILE", CHINESE_FILE.files[0]);
    formData.append("ORIGINAL_FILE", ORIGINAL_FILE.files[0]);
    formData.append("REVISED_FILE", REVISED_FILE.files[0]);

    fetch('/comparefiles', {
        method: "POST", 
        body: formData
    })
    .then(response => response.blob())
    .then(blob => {
        console.log("size: " + blob.size)
        download(blob, "output.xlsx")
    });
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