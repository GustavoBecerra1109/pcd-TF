fetch('/getMessages')
  .then(response => response.json())
  .then(data => {
    // data contiene los mensajes en formato JSON
    // Haz lo que necesites con los mensajes
    console.log(data);
  })
  .catch(error => {
    console.error('Error:', error);
  });