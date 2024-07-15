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
    let form = document.getElementById("formRenew");

    // disable form elements
    for (let i = 0; i < form.elements.length; i++) {
        $(form.elements[i]).attr("disabled", true);
        if ($(form.elements[i]).attr("type") == "submit") {
            $(form.elements[i]).find("i").
                removeClass("fa-solid fa-floppy-disk").
                addClass("spinner-border spinner-border-sm");
        }
    }

    $('#btn-logout').attr('disabled', true);
}

function ajaxDone() {
    let form = document.getElementById("formRenew");

    // disable form login
    for (let i = 0; i < form.elements.length; i++) {
        $(form.elements[i]).removeAttr("disabled");
        if ($(form.elements[i]).attr("type") == "submit") {
            $(form.elements[i]).find("i").
                removeClass("spinner-border spinner-border-sm").
                addClass("fa-solid fa-floppy-disk");
        }
    }

    $('#btn-logout').removeAttr('disabled');
}

function renewPassword(form) {
    $.ajax({
        type: $(form).attr('method'),
        url: $(form).attr('action'),
        cache: false,
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Accept', '*/*');
            xhr.setRequestHeader('Authorization', document.getElementById('tokenBox').innerHTML);
            ajaxLoading();
        },
        data: $(form).serialize(),
        dataType: "json"
    }).done(function (response) {
        // go to submit token page
        window.location.href = document.getElementById('redirectPath').innerHTML
    }).fail(function ($jqXHR, $errorThrown) {
        ajaxDone();
        try {
            let response = JSON.parse($jqXHR.responseText);
            if (response.code) {
                switch (response.code) {
                    case 400003:
                        toastInfo("anda tidak dapat menggunakan password yang sama!");
                        break;
                    case 400999:
                        // validator fails
                        for (let x in response.message) {
                            if (x == "password") {
                                for (let y in response.message[x]) {
                                    if (response.message[x][y].tag == "required") {
                                        toastInfo("password wajib di isi!");
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