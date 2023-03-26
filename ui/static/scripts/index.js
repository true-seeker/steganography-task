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
            });
        } else {
            response = await fetch('api/pic_to_pic', {
                method: 'POST',
                body: formData,
            });
        }
        let result = await response;
        console.log(result)

        return false
    }
    form.classList.add('was-validated')
}, false)