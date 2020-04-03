firebase.initializeApp({
    projectId: "codebender-12e17",
    databaseURL: "https://codebender-12e17.firebaseio.com"
});
const fs = firebase.firestore();

const authorddl = document.querySelector('#author-dropdown');
const form = document.querySelector('#new-message-form');
const conversation = document.querySelector('#message-list');

form.addEventListener('submit', (e) => {
    e.preventDefault();
    const request = {
        conversation_id: "123",
        author: authorddl.value,
        message: form.message.value
    }
    console.log(request)

    // Send request to GCF
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "https://us-central1-codebender-12e17.cloudfunctions.net/SendMessage", true);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

    // send the collected data as JSON
    xhr.send(JSON.stringify(request));

    xhr.onloadend = function () {
        form.message.value = '';
    };
});

// render chat message
function renderMessage(doc){
    let li = document.createElement('li');
    let message = document.createElement('span');
    let author = document.createElement('span');
    let timestamp = document.createElement('span');

    li.setAttribute('data-id', doc.id);
    author.textContent = doc.data().author;
    message.textContent = doc.data().message;

    var date = new Date(0);
    date.setUTCSeconds(doc.data().timestamp);
    timestamp.textContent = date;

    li.appendChild(author);
    li.appendChild(message);
    li.appendChild(timestamp);
    conversation.appendChild(li);
}


// Firestore listener
fs.collection('conversations/123/messages').orderBy('timestamp').onSnapshot(snapshot => {
    let changes = snapshot.docChanges();
    changes.forEach(change => {
        console.log(change.doc.data());
        renderMessage(change.doc);
    });
});