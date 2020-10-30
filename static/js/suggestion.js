function setSuggestions(input, picker, options) {
    input.onkeyup = () => {
        clearTimeout(timer)
        timer = setTimeout(() => {
            options.query = input.value
            fetch('/suggestion', {
                method: 'POST',
                mode: 'cors',
                body: JSON.stringify(options)
            }).then(response => response.text()).then(result => {
                picker.innerHTML = ""
                JSON.parse(result).suggestions.forEach(e => {
                    let line = document.createElement('div')
                    line.onclick = () => {
                        picker.innerHTML = ""
                        input.value = e.value.includes(cityPicker.value) ? e.value.replace(cityPicker.value + ', ', '') : e.value
                    }
                    line.innerText = e.value
                    picker.appendChild(line)
                })

            }).catch(error => console.log("error", error))
        }, threshhold * 500)
    }
}

const cityPicker = document.getElementById("city")
const cities = document.getElementById("cityPicker")

setSuggestions(cityPicker, cities, 
})

const housePicker = document.getElementById("house")
const houses = document.getElementById("housePicker")

setSuggestions(housePicker, houses, {
    "from_bound": {
        "value": "street"
    },
    "to_bound": {
        "value": "house"
    },
    "locations": [{
        "region": () => {
            return cityPicker.value
        }
    }],
})