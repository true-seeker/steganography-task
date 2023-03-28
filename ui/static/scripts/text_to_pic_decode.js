const form = document.querySelector('#encrypt-form')
const submit_button = document.querySelector('#submit-button')
submit_button.addEventListener('click', async event => {
        event.preventDefault()
        event.stopPropagation()
        if (!form.checkValidity()) {
            console.log("false valid")
        } else {
            let formData = new FormData(form);
            await fetch('/api/text_to_pic_decode', {
                    method: 'POST',
                    body: formData,
                }
            ).then(response => response.json())
                .then(response => {
                    document.getElementById('message-modal-text').innerHTML = response.message;

                    const myModal = new bootstrap.Modal(document.getElementById('message-modal'), {})
                    myModal.show()
                })
            return
        }
        form.classList.add('was-validated')
    }
    ,
    false
)