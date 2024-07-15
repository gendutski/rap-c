function toastInfo(message) {
    Toastify({
        text: message,
        duration: 3000,
        close: true,
        gravity: "top",
        position: "right",
        style: {
            background: "#2c1c19"
        }
    }).showToast();
}

function ajaxLoading() {
    let form = document.getElementById("formReset");

    // disable form elements
    for (let i = 0; i < form.elements.length; i++) {
        $(form.elements[i]).attr("disabled", true);
        if ($(form.elements[i]).attr("type") == "submit") {
            $(form.elements[i]).find("i").
                removeClass("fa-solid fa-floppy-disk").
                addClass("spinner-border spinner-border-sm");
        }
    }
}

function ajaxDone() {
    let form = document.getElementById("formReset");

    // disable form login
    for (let i = 0; i < form.elements.length; i++) {
        $(form.elements[i]).removeAttr("disabled");
        if ($(form.elements[i]).attr("type") == "submit") {
            $(form.elements[i]).find("i").
                removeClass("spinner-border spinner-border-sm").
                addClass("fa-solid fa-floppy-disk");
        }
    }
}

function resetPassword(form) {
    $.ajax({
        type: $(form).attr('method'),
        url: $(form).attr('action'),
        cache: false,
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Accept', '*/*');
            ajaxLoading();
        },
        data: $(form).serialize(),
        dataType: "json"
    }).done(function (response) {
        toastInfo("password tersimpan");

        // go to submit token page
        setTimeout(function () {
            $('#formSubmitToken input[name="token"]').val(response.token);
            $('#formSubmitToken').submit();
        }, 1500);
    }).fail(function ($jqXHR) {
        ajaxDone();
        try {
            let response = JSON.parse($jqXHR.responseText);
            if (response.code) {
                switch (response.code) {
                    case 404001:
                        toastInfo("token tidak ditemukan, atau token sudah kadaluwarsa");
                        break;
                    case 404002:
                        toastInfo("email tidak ditemukan");
                        break;
                    case 400999:
                        // validator fails
                        for (let x in response.message) {
                            if (x == "token") {
                                for (let y in response.message[x]) {
                                    if (response.message[x][y].tag == "required") {
                                        toastInfo("token dari email wajib disertakan!");
                                    }
                                }
                            } else if (x == "password") {
                                for (let y in response.message[x]) {
                                    if (response.message[x][y].tag == "required") {
                                        toastInfo("password wajib diisi!");
                                    }
                                }
                            } else if (x == "confirmPassword") {
                                for (let y in response.message[x]) {
                                    if (response.message[x][y].tag == "required") {
                                        toastInfo("konfirmasi password wajib diisi!");
                                    } else if (response.message[x][y].tag == "eqfield") {
                                        toastInfo("konfirmasi password tidak sama!");
                                    }
                                }
                            } else {
                                toastInfo("data yang diinput tidak valid");
                            }
                        }
                        break
                    default:
                        toastInfo("ada kesalahan teknis, error #" + response.code);
                        break;
                }
            }
        } catch (error) {
            toastInfo("ada kesalahan teknis");
            console.log(error);
        }
        return;
    });
}