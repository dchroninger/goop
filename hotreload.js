var socket = new WebSocket(`ws://localhost:5555/`);

socket.onopen = function () {
  console.info("Hot reload WebSocket connected");
};

socket.onclose = function () {
  //TODO: Find a better way to do this
  setTimeout(() => {
    window.location.reload();
  }, 800);
};
