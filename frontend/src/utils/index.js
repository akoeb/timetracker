// utility functions:
function createProject(name, clientName, status) {
  const data = {
    name: name,
    client_name: clientName,
    status: status
  };
  let retval = {};
  fetch("http://localhost:8080/api/projects", {
    method: "POST", // or 'PUT'
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(data)
  })
    .then(response => response.json())
    .then(data => {
      console.log("Success:", data);
      retval = data;
    })
    .catch(error => {
      console.error("Error:", error);
    });
  return retval;
}

function deleteProject(id) {
  let retval = {};
  fetch("http://localhost:8080/api/projects/" + id, {
    method: "DELETE", // or 'PUT'
    headers: {
      "Content-Type": "application/json"
    }
  })
    .then(response => response.json())
    .then(data => {
      console.log("Success:", data);
      retval = data;
    })
    .catch(error => {
      console.error("Error:", error);
    });
  return retval;
}

export { createProject, deleteProject };
