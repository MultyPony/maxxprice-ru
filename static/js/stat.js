try {
    document.addEventListener("DOMContentLoaded", function () {
        // var loader = new ldBar("#loader", {
        //     "preset": "energy",
        //     "value": 0
        // });
        var loader = new ldLoader({ root: "#loader" }); 

        IMask(document.getElementById('phone'), {
            mask: '+{7}(000)000-00-00'
        })
        IMask(document.getElementById('phone_rec'), {
            mask: '+{7}(000)000-00-00'
        })
        IMask(document.getElementById('phone_submit'), {
            mask: '+{7}(000)000-00-00'
        })
        const result_cont = document.getElementsByClassName('calculateresult')[0]

        let mainButton = document.getElementById('getav');

        let city = "",
            house = "",
            square = "",
            stair = "",
            rem = ""
        document.cityTail = new SlimSelect({
            select: document.getElementById("city"),
            placeholder: 'Город',
            searchPlaceholder: "Выберите из списка",
            onChange: (item) => {
                city = item.value != "Город" ? item.value : ""
                document.getElementById('house').classList.remove('disabled-input')
            }
        })
        document.houseTail = new SlimSelect({
            select: document.getElementById("house"),
            placeholder: 'Улица, Дом',
            searchPlaceholder: "Введите свой адрес и выберите его в появившемся ниже списке",
            search: true,
            searchFilter: (option, search) => {
                let found = false
                search.split(' ').forEach((e) => {
                    // console.log(option.text)
                    // console.log(e)
                    if (e.length > 1 || (e.length == 1 && e * 1 != NaN))
                        if (option.text.toLowerCase().includes(e.toLowerCase())) {
                            found = true
                        }
                })
                return found
            },
            ajax: (search, callback) => {
                fetch('/suggestion', {
                    method: 'POST',
                    mode: 'cors',
                    body: JSON.stringify({
                        "query": search,
                        "locations": [{
                            "region": "москва"
                        }, {
                            "region": "московская обл"
                        }, ],
                        "from_bound": {
                            "value": "street"
                        },
                        "to_bound": {
                            "value": "house"
                        }
                    })
                }).then(response => response.text()).then(result => {
                    function cntContains(haystack, needles) {
                        let cnt = 0

                        needles.forEach((e) => {
                            if (haystack.includes(e)) {
                                cnt++
                            }
                        })

                        return cnt
                    }
                    let resp = JSON.parse(result).suggestions
                    // if (resp.length > 1) {
                    //     resp = resp.sort((a, b) => {
                    //         return cntContains(a.value, search.split(' ')) > cntContains(b.value, search.split(' '))
                    //     })
                    // }
                    document.getElementById('square').classList.remove('disabled-input')
                    let res = []
                    resp.forEach(e => {
                        //console.log(e.value)
                        res.push({
                            text: e.value,
                            value: e.value,
                        })
                    })
                    callback(res)
                }).catch(error => console.log("error", error))
            },
        });
        document.querySelector('#house+div > .ss-single-selected').onclick = () => {
            document.houseTail.open()
            setTimeout(() => {
                document.querySelector('#house+div .ss-search input').value = house
            }, 100)
        }
        document.squareTail = new SlimSelect({
            select: document.getElementById("square"),
            placeholder: 'Жил. площадь',
            searchPlaceholder: "Выберите из списка",
            onChange: (item) => {
                square = item.value != "Жил. площадь" ? item.value : ""
                document.getElementById('stair').classList.remove('disabled-input')
            }
        });
        document.stairTail = new SlimSelect({
            select: document.getElementById("stair"),
            searchPlaceholder: "Выберите из списка",
            placeholder: 'Этаж',
            onChange: (item) => {
                stair = item.value != "Этаж" ? item.value : ""
                document.getElementById('rem').classList.remove('disabled-input')
            }
        })
        document.remTail = new SlimSelect({
            select: document.getElementById("rem"),
            searchPlaceholder: "Выберите из списка",
            placeholder: 'Отделка',
            onChange: (item) => {
                rem = item.value != "Отделка" ? item.value : ""
                if (city != "" && house != "" && square != "" && stair != "" && rem != "") {
                    mainButton.classList.remove('disabled')
                    // document.getElementById('getav').classList.remove('disabled')
                }
            }
        })

        // document.getElementById('getav').addEventListener("click", loadButton)
        mainButton.addEventListener("click", loadButton)
        
        function loadButton() {
            let fixed = document.getElementsByClassName("loader-wrap")
            fixed[0].classList.remove("none");
            loader.on();
            setTimeout(() => {
              loader.off();
              fixed[0].classList.add('none');
              document.getElementsByClassName("submit-contacts")[0].classList.remove("none");
            }, 2000);
            mainButton.classList.add('disabled')
          }

        let submitContacts = document.getElementById("submit-contacts");
        submitContacts.addEventListener("click", async function(event) {
            event.preventDefault();
            let response = await fetch("/zhopa", {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json;charset=utf-8'
                },
                body: JSON.stringify("kekes")
              });
            if (response.ok) { // если HTTP-статус в диапазоне 200-299
                // получаем тело ответа (см. про этот метод ниже)
                let json = await response.json();
              } else {
                alert("Ошибка HTTP: " + response.status);
              }
            console.log(response);
        })
        // document.getElementById('getav').onclick = () => {
        //     if (document.houseTail.selected() != "") {
        //         document.getElementById('getav').classList.remove('disabled')
        //     }
        //     if (document.getElementById('getav').classList.contains('disabled')) {
        //         if (city == "") {
        //             document.getElementById('city').nextElementSibling.classList.add('err')
        //             setTimeout(() => {
        //                 document.getElementById('city').nextElementSibling.classList.remove('err')
        //             }, 2000)
        //         }
        //         if (house == "") {
        //             document.getElementById('house').nextElementSibling.classList.add('err')
        //             setTimeout(() => {
        //                 document.getElementById('house').nextElementSibling.classList.remove('err')
        //             }, 2000)
        //         }
        //         if (square == "") {
        //             document.getElementById('square').nextElementSibling.classList.add('err')
        //             setTimeout(() => {
        //                 document.getElementById('square').nextElementSibling.classList.remove('err')
        //             }, 2000)
        //         }
        //         if (stair == "") {
        //             document.getElementById('stair').nextElementSibling.classList.add('err')
        //             setTimeout(() => {
        //                 document.getElementById('stair').nextElementSibling.classList.remove('err')
        //             }, 2000)
        //         }
        //         if (rem == "") {
        //             document.getElementById('rem').nextElementSibling.classList.add('err')
        //             setTimeout(() => {
        //                 document.getElementById('rem').nextElementSibling.classList.remove('err')
        //             }, 2000)
        //         }
        //     } else {

        //         document.getElementById('loader').classList.remove('hidden')
        //         document.getElementById('loader').previousElementSibling.classList.remove('hidden')

        //         fetch('/api/categories', {
        //             method: "POST",
        //             body: JSON.stringify({
        //                 rem: rem * 1,
        //                 house: document.houseTail.selected(),
        //                 square: square.split(' ')[0] * 1,
        //                 city: city,
        //                 stair: stair,
        //             })
        //         }).then(response => response.json()).then(average => {
        //             house = document.houseTail.selected()
        //             let ind = 100
        //             var p = document.querySelectorAll('.ldBar-label')[0];
        //             var c = 59
        //             setInterval(() => {
        //                 p.setAttribute('data-value', "00:" + c)
        //                 if (c > 0) c--
        //             }, 1000)
        //             let timerid = setInterval(() => {
        //                 loader.set(--ind, false)
        //                 if (ind % 10 == 0) {
        //                     fetch('/api/average', {
        //                         method: "POST",
        //                         body: JSON.stringify({
        //                             rem: average.avg,
        //                             house: document.houseTail.selected(),
        //                             square: square.split(' ')[0] * 1,
        //                             city: city
        //                         })
        //                     }).then(response => response.json()).then(average => {
        //                         if (average.avg == 0) {
        //                             if (ind == 100)
        //                                 document.getElementById('full').classList.add('shown')
        //                             clearInterval(timerid)

        //                         } else {
        //                             clearInterval(timerid)

        //                             document.getElementById('loader').classList.add('hidden')
        //                             document.getElementById('loader').previousElementSibling.classList.add('hidden')

        //                             result_cont.querySelectorAll('p').forEach((e) => e.remove())
        //                             let fsttextbloc = document.createElement('p')
        //                             let sndtextbloc = document.createElement('p')
        //                             fsttextbloc.innerText = `Цена объекта от ` + average.avg + ` до ` + average.agv + ` руб.`
        //                             sndtextbloc.innerText = `Отправьте нам заявку, мы соберём предложения всех компаний выкупающих объекты недвижимости, предложим больше всех, чтобы вы не тратили своё время на поиски`
        //                             result_cont.prepend(sndtextbloc)
        //                             result_cont.prepend(fsttextbloc)
        //                             result_cont.classList.remove('hidden')
        //                         }
        //                     })
        //                 }
        //             }, 600)
        //         })

        //     }

        // }
        document.getElementById('phone').onkeyup = () => {
            document.getElementById('sendRecall').classList.remove('disabled')
        }
        const phoneInput = document.getElementById('phone')
        document.getElementById('sendRecall').onclick = () => {
            fetch('/api/recall', {
                method: "POST",
                body: JSON.stringify({
                    phone: phoneInput.value,
                    rem: rem,
                    stair: stair,
                    city: city,
                    house: document.houseTail.selected(),
                    square: square
                })
            })
            alert('Ваша заявка напралена на обработку')
            document.getElementById('rec').classList.remove('shown')

        }
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
                    stair: stair,
                    city: city,
                    house: house,
                    square: square
                })
            })
            alert('Мы вам перезвоним!')
            document.getElementById('rec').classList.remove('shown')
        }
        document.getElementById('sendres').onclick = () => {
            let phone = document.getElementById('recrec').value
            fetch('/api/recall', {
                method: "POST",
                body: JSON.stringify({
                    phone: phone,
                    rem: rem,
                    stair: stair,
                    city: city,
                    house: house,
                    square: square
                })
            })
            alert('Ожидайте результат по указанному адресу')
            document.getElementById('full').classList.remove('shown')
        }
        document.getElementById('rec').getElementsByClassName('modal-close')[0].onclick = () => {
            document.getElementById('rec').classList.remove('shown')
        }
        document.getElementById('full').getElementsByClassName('modal-close')[0].onclick = () => {
            document.getElementById('full').classList.remove('shown')
        }
    });

} catch (error) {
    if (location.search == "?q=test") {
        alert(error)
    }
}