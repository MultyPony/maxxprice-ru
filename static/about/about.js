document.getElementById("phone_rec").onkeyup = (value) => {
    if (document.getElementById("phone_rec").value.length == 16) {
        document.getElementById('sendrec').classList.remove('disabled')
    } else {
        document.getElementById('sendrec').classList.add('disabled')
    }
}
document.getElementById('tr').onclick = () => {
    document.getElementById('rec').classList.add('shown')
}
document.getElementById('sendrec').onclick = () => {
    let phone = document.getElementById('phone_rec').value
    fetch('/api/recall', {
        method: "POST",
        body: JSON.stringify({
            phone: phone,
            rem: rem,
            rooms: rooms,
            city: city,
            house: house,
            square: square,
        })
    })
    document.getElementById("phone_rec").value = ''
    document.getElementById('sendrec').classList.add('disabled')
    // alert('Мы вам перезвоним!')
    document.getElementById('rec').classList.remove('shown')
}