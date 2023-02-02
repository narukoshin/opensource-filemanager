let download = (element) => {
    let file_name = element.getAttribute("download-id")
    window.open("/download/" + file_name)
}

let openFolder = (element) => {
    let folder_name = element.getAttribute("folder-id")
    let file_list   = document.querySelector(".file-list")
    fetch("/folder/?folder_name=" + folder_name)
    .then(response => response.text())
    .then(data => {
        file_list.innerHTML = data
    })
}

let resize_file_names = () => {
    let w = window.innerWidth
    let file_name = document.getElementsByClassName("file-name")
    let file_name_new_length = Math.floor((w / 100) + 5)
    if (w < 1000) {
        Array.from(file_name).forEach((element) => {
          if ((element.textContent.length > file_name_new_length) && (element.textContent.length != file_name_new_length)) {
            element.textContent = element.textContent.substr(0, file_name_new_length) + "..."
          }  
        })
    }
}
resize_file_names()