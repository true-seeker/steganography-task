export function get_file_name_from_response(result) {
    if (!result.ok) {
        console.log(result.statusText);
    }

    const header = result.headers.get('Content-Disposition');
    const parts = header.split(';');
    return parts[1].split('=')[1].replaceAll("\"", "")
}

export function prepare_file(result) {
    if (!result.ok) {
        console.log(result.statusText);
    }
    return result.blob();
}

export function save_file(blob, filename) {
    if (blob != null) {
        let url = window.URL.createObjectURL(blob);
        let a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        a.remove();
    }
}

export function show_toast(color, header, message, toast_id) {
    document.getElementById('toast_rect').setAttribute('fill', color);
    document.getElementById('toast_header').innerHTML = header;
    document.getElementById('toast_message').innerHTML = message;
    const toast = document.getElementById(toast_id)
    console.log(toast)
    const toastBootstrap = bootstrap.Toast.getOrCreateInstance(toast)
    toastBootstrap.show()
}