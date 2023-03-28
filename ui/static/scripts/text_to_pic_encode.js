import {get_file_name_from_response, prepare_file, save_file, show_toast} from "./utils.js";

const form = document.querySelector('#encrypt-form')
const submit_button = document.querySelector('#submit-button')
submit_button.addEventListener('click', async event => {
        event.preventDefault()
        event.stopPropagation()
        if (!form.checkValidity()) {
            console.log("false valid")
        } else {
            let formData = new FormData(form);
            await fetch('/api/text_to_pic_encode', {
                    method: 'POST',
                    body: formData,
                }
            ).then((response) => {
                    let filename = get_file_name_from_response(response)
                    prepare_file(response)
                        .then((blob) => save_file(blob, filename))
                    show_toast('#95CD41', 'Success', 'Encoded picture download will start soon', 'alert-toast');
                }
            )
            return
        }
        form.classList.add('was-validated')
    }
    ,
    false
)