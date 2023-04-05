window.e_receiveConsoleOutput = function (output) {
  const consoleOutput = document.getElementById("consoleOutput");

  consoleOutput.innerHTML += output + "\n";
  consoleOutput.scrollTop = consoleOutput.scrollHeight;
};

window.e_receiveRunningState = function (running) {
  console.log("Running status: " + running);
  const consoleStatus = document.getElementById("consoleStatus");
  consoleStatus.classList = running ? "running" : "stopped";

  const forceStop = document.getElementById("forceStop");

  if (running) {
    disableButtons();
    forceStop.disabled = false;
  } else {
    enableButtons();
    forceStop.disabled = true;
  }
};

function disableButtons() {
  const buttons = document.querySelectorAll("#buttons button");
  for (const button of buttons) {
    button.disabled = true;
  }
}

function enableButtons() {
  const buttons = document.querySelectorAll("#buttons button");
  for (const button of buttons) {
    button.disabled = false;
  }

  if (document.getElementById("gameSelectorSelection").options.length > 1) {
    document.getElementById("stream").disabled = false;
  } else {
    document.getElementById("stream").disabled = true;
  }
}

async function getGames() {
  const gameSelectorSelection = document.getElementById(
    "gameSelectorSelection"
  );

  const unknownGame = document.getElementById("gameSelectorUnknownOption");
  unknownGame.selected = false;
  unknownGame.hidden = true;

  function doRemove() {
    const toRemove = [];
    for (const option of gameSelectorSelection.options) {
      if (option === unknownGame) {
        continue;
      }

      toRemove.push(option);
    }

    for (const option of toRemove) {
      option.remove();
    }
  }

  window
    .b_getGames()
    .then((games) => {
      if (games.length === 0) {
        return;
      }

      doRemove();

      games.sort();
      for (const game of games) {
        const option = document.createElement("option");
        option.value = game;
        option.innerText = game;

        gameSelectorSelection.appendChild(option);
      }

      gameSelectorSelection.disabled = false;
      gameSelectorSelection.value = games[0];

      document.getElementById("stream").disabled = false;
    })
    .catch((err) => {
      doRemove();

      gameSelectorSelection.disabled = true;
      unknownGame.hidden = false;
      unknownGame.selected = true;

      console.error(err);

      document.getElementById("stream").disabled = true;
    });
}

async function streamGame() {
  const gameSelectorSelection = document.getElementById(
    "gameSelectorSelection"
  );
  const game = gameSelectorSelection.value;

  const mode = document.getElementById("mode").value;
  const bitrate = document.getElementById("bitrate").value;
  const fps = document.getElementById("fps").value;
  const resolution = document.getElementById("resolution").value;

  window
    .b_streamGame({
      bitrate: Number(bitrate) * 1000,
      fps: Number(fps),
      game: game,
      mode: mode === "net" ? "dji_net" : "dji_usb",
      resolution: {
        width: Number(resolution.split("x")[0]),
        height: Number(resolution.split("x")[1]),
      },
    })
    .then(() => {})
    .catch((err) => {
      console.error(err);
    });
}

function doBinds() {
  const pair = document.getElementById("pair");
  pair.addEventListener("click", () => {
    window
      .b_pair()
      .then(() => {})
      .catch((err) => {
        console.error(err);
      });
  });

  const unpair = document.getElementById("unpair");
  unpair.addEventListener("click", () => {
    window
      .b_unpair()
      .then(() => {})
      .catch((err) => {
        console.error(err);
      });
  });

  const refresh = document.getElementById("gameSelectorRefresh");
  refresh.addEventListener("click", getGames);

  const forceStop = document.getElementById("forceStop");
  forceStop.addEventListener("click", () => {
    window
      .b_forceStop()
      .then(() => {})
      .catch((err) => {
        console.error(err);
      });
  });

  const stream = document.getElementById("stream");
  stream.addEventListener("click", streamGame);

  const quit = document.getElementById("quit");
  quit.addEventListener("click", () => {
    window
      .b_quit()
      .then(() => {})
      .catch((err) => {
        console.error(err);
      });
  });
}

window.addEventListener("DOMContentLoaded", () => {
  doBinds();
});
