function show_text() {
    document.getElementById('file-input-div').hidden = true;
    document.getElementById('source-file').required = false;
    document.getElementById('source-file').value = null;
    document.getElementById('text-input-div').hidden = false;
    document.getElementById('input-text').required = true;
}

function show_file_input() {
    document.getElementById('text-input-div').hidden = true;
    document.getElementById('input-text').required = false;
    document.getElementById('input-text').value = null;
    document.getElementById('file-input-div').hidden = false;
    document.getElementById('source-file').required = true;
}

function prepare_file(result) {
    if (!result.ok) {
        throw Error(result.statusText);
    }

    // We are reading the *Content-Disposition* header for getting the original filename given from the server
    const header = result.headers.get('Content-Disposition');
    const parts = header.split(';');
    filename = parts[1].split('=')[1].replaceAll("\"", "");
    return result.blob();
}

function save_file(blob) {
    if (blob != null) {
        var url = window.URL.createObjectURL(blob);
        var a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        a.remove();
    }
}

const form = document.querySelector('#encrypt-form')

const submit_button = document.querySelector('#submit-button')
submit_button.addEventListener('click', async event => {
        event.preventDefault()
        event.stopPropagation()
        if (!form.checkValidity()) {
            console.log("false valid")
        } else {
            let formData = new FormData(form);

            let response;
            if (formData.get('stegoType') === 'text') {
                response = await fetch('api/text_to_pic', {
                        method: 'POST',
                        body: formData,
                    }
                ).then((result) => prepare_file(result).then((blob) => save_file(blob)))
            } else {
                response = await fetch('api/pic_to_pic', {
                    method: 'POST',
                    body: formData,
                }).then((result) => prepare_file(result).then((blob) => save_file(blob)));
            }
            let result = await response;
            console.log(result)

            return false
        }
        form.classList.add('was-validated')
    }
    ,
    false
)