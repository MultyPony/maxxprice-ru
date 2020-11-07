var loader = new ldLoader({ root: "#loader-feedback" }); 

document.getElementById('load-more').addEventListener("click", loadButton)
        
function loadButton() {
    let fixed = document.getElementsByClassName("loader-wrap")
    let errorText = document.getElementsByClassName("error-feedback")[0]
    errorText.classList.add('none');
    fixed[0].classList.remove("none");
    loader.on();
    setTimeout(() => {
        loader.off();
        fixed[0].classList.add('none');
        errorText.classList.remove('none');
    }, 2000);
    scrollingElement = document.getElementById('feedback');
    scrollingElement.scrollTop = scrollingElement.scrollHeight;
}

let leaveFeedbackForm = document.getElementById('leave-feedback-form');

document.getElementById('leave-feedback-btn').onclick = () => {
    leaveFeedbackForm.classList.add('shown')
}
leaveFeedbackForm.getElementsByClassName('modal-close')[0].onclick = () => {
    leaveFeedbackForm.classList.remove('shown')
}

let leaveFeedbackName = document.getElementById('leave-feedback-name');
let leaveFeedbackLotNumber = document.getElementById('leave-feedback-lotnumber');
let feedbackTextarea = document.getElementById('feedback-textarea');
let leaveFeedbackSubmit = document.getElementById('leave-feedback-submit');


leaveFeedbackName.addEventListener("keyup", function() {
    // console.log('feedbackTextarea.value:' + feedbackTextarea.value);
    if (feedbackTextarea.value !== '' && leaveFeedbackName.value !== '' && leaveFeedbackLotNumber.value !== '') {
        leaveFeedbackSubmit.classList.remove('disabled');
    } 
    else {
        leaveFeedbackSubmit.classList.add('disabled');
    }
});

leaveFeedbackLotNumber.addEventListener("keyup", function() {
    // console.log('feedbackTextarea.value:' + feedbackTextarea.value);
    if (feedbackTextarea.value !== '' && leaveFeedbackName.value !== '' && leaveFeedbackLotNumber.value !== '') {
        leaveFeedbackSubmit.classList.remove('disabled');
    } 
    else {
        leaveFeedbackSubmit.classList.add('disabled');
    }
});

feedbackTextarea.addEventListener("keyup", function() {
    console.log('feedbackTextarea.value:' + feedbackTextarea.value);
    if (feedbackTextarea.value !== '' && leaveFeedbackName.value !== '' && leaveFeedbackLotNumber.value !== '') {
        leaveFeedbackSubmit.classList.remove('disabled');
    } 
    else {
        leaveFeedbackSubmit.classList.add('disabled');
    }
});


document.getElementById('leave-feedback-submit').onclick = () => {
    if (leaveFeedbackName.value === '') {
        leaveFeedbackSubmit.classList.add('disabled');
        return

    }
    else if (leaveFeedbackLotNumber.value === '') {
        leaveFeedbackSubmit.classList.add('disabled');
        return

    }
    else if (feedbackTextarea.value === '') {
        leaveFeedbackSubmit.classList.add('disabled');
        return

    }
    else {
        leaveFeedbackSubmit.classList.remove('disabled');
        leaveFeedbackName.value = '';
        leaveFeedbackLotNumber.value = '';
        feedbackTextarea.value = '';
        leaveFeedbackSubmit.value = '';
        leaveFeedbackForm.classList.remove('shown');
    }
}