// utility functions:
/*
class Project {
  constructor(name, clientName, status) {
    this.id = 0;
    this.name = name;
    this.clientName = clientName;
    this.status = status;
  }
}
*/
window.onload = function() {
    listProjects()
}
function listProjects() {
  fetch("http://localhost:8080/api/projects", {
    method: "GET", // or 'PUT'
    headers: {
      "Content-Type": "application/json"
    }
  })
    .then(response => {
      if (response.status != 200) {
        throw Error("Unexpected Status Code: " + response.statusText);
      }
      return response.json()
    })
    .then(data => {
      displayProjectList(data)
    })
    .catch(error => {
      console.error("Error:", error);
    });
}
function displayProjectList(projectList) {
    // get mount point in page:
    let app = document.getElementById('app')
  
    // get templates
    let projectLustWrapperTpl = document.getElementById('project-overview-tpl')
    let projectListTpl = document.getElementById('project-list-tpl')
    let projectItemTpl = document.getElementById('project-item-tpl')
  
    // clone wrapper element templates
    var list_clone = document.importNode(projectListTpl.content, true)
    var list_wrapper_clone = document.importNode(projectLustWrapperTpl.content, true)

    // early exit if we do not have projects to display:
    if (!projectList.projects || projectList.projects.length < 1) {
      list_wrapper_clone.getElementById('project-list').textContent = "No projects to display"
    }
    else {
        // fill project item template and add it to project list template
        for (project of projectList.projects) {
            // clone tpl
            var item_clone = document.importNode(projectItemTpl.content, true);

            // fill with values
            item_clone.querySelector('.project-name-val').value = project.name
            item_clone.querySelector('.project-client-name-val').value = project.client_name

            // add status select dropdown:
            let sel = statusSelect(project.status).getElementById('status-select')
            sel.name = "ProjectStatus"
            sel.removeAttribute('id')
            sel.disabled = true
            item_clone.querySelector('.project-status-val').appendChild(sel)

            // add identifier to the buttons
            item_clone.querySelector('.project-show-btn').dataset.projectid = project.id
            item_clone.querySelector('.project-del-btn').dataset.projectid = project.id
            item_clone.querySelector('.project-edit-btn').dataset.projectid = project.id

            // append to wrapper element
            list_clone.getElementById('project-item').appendChild(item_clone)
        }

        // project list to project list_wrapper
        list_wrapper_clone.getElementById('project-list').appendChild(list_clone)
  }

  // now activate the create project form
  var sel = statusSelect().getElementById('status-select')
  sel.name = "newProjectStatus"
  form = document.getElementById('add-project-form')
  form.insertBefore(sel, document.getElementById('add-project-form-submit'))
  form.addEventListener('submit', addProject);

  // and display everything in the page
  app.textContent = ''
  app.appendChild(overview_clone)
}

function statusSelect(status) {
  const tpl = document.getElementById('status-select-tpl');
  let clone = document.importNode(tpl.content, true);
  if (status) {
    clone.getElementById('status-select').value = status;
  }
  return clone
}


function addProject(ev) {
  ev.preventDefault()
  const form=document.getElementById('add-project-form')
  createProject(form.newProjectName.value, form.newProjectClientName.value, form.newProjectStatus.value)

}
function createProject(name, clientName, status) {
  const data = {
    name: name,
    client_name: clientName,
    status: status
  };
  fetch("http://localhost:8080/api/projects", {
    method: "POST", // or 'PUT'
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(data)
  })
    .then(response => {
      if (response.status != 201) {
        throw Error("Unexpected Status Code: " + response.statusText);
      }
      return response.json()
    })
    .then(data => {
      listProjects();
    })
    .catch(error => {
      console.error("Error:", error);
    });
}

function toggleEditProject(elem) {
    console.log("elem" + elem)
    // get projectid from calling element
    projectId = elem.dataset.projectid
    
    // detect edit vs update mode from button content
    if (elem.textContent == "update") {
      
      updateProject(
        projectId,
        elem.parentNode.querySelector('.project-name-val').value,
        elem.parentNode.querySelector('.project-client-name-val').value,
        elem.parentNode.querySelector('.status-select').value
      );

      elem.textContent = "edit"
      elem.parentNode.querySelector('.project-name-val').disabled = true
      elem.parentNode.querySelector('.project-client-name-val').disabled = true
      elem.parentNode.querySelector('.status-select').disabled = true    
    }
    else if (elem.textContent == "edit"){
      elem.textContent = "update"
      elem.parentNode.querySelector('.project-name-val').disabled = false
      elem.parentNode.querySelector('.project-client-name-val').disabled = false
      elem.parentNode.querySelector('.status-select').disabled = false    
    }
    else {
      console.log("wrong button value")
    }
}

function updateProject(id, name, clientName, status) {
  // get projectid from calling element
  const data = {
    id: parseInt(id),
    name: name,
    client_name: clientName,
    status: status
  };
  console.log(data)
  fetch("http://localhost:8080/api/projects/" + id, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(data)
  })
  .then(response => {
      console.log("Success:", response);
      if (response.status != 200) {
        throw Error("Unexpected Status Code: " + res.statusText);
      }
      listProjects();
  })
  .catch(error => {
    console.error("Error:", error);
  });
}

function deleteProject(elem) {
    // get projectid from calling element
    projectId = elem.dataset.projectid
    fetch("http://localhost:8080/api/projects/" + projectId, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json"
      }
    })
    .then(response => {
        console.log("Success:", response);
        if (response.status != 204) {
          throw Error("Unexpected Status Code: " + res.statusText);
        }
        listProjects();
    })
    .catch(error => {
      console.error("Error:", error);
    });
}

function showProject(elem) {
    // get projectid from calling element
    projectId = elem.dataset.projectid
    fetch("http://localhost:8080/api/projects/" + projectId, {
        method: "GET", // or 'PUT'
        headers: {
            "Content-Type": "application/json"
        }
     })
    .then(response => {
      if (response.status != 200) {
        throw Error("Unexpected Status Code: " + response.statusText);
      }
      return response.json()
    })
    .then(data => {
      displayProject(data)
    })
    .catch(error => {
      console.error("Error:", error);
    });
}

function displayProjectHistory(project) {
  // get mount point in page:
  let app = document.getElementById('app')

  // prepare item-template:
  let projectEventItemTpl = document.getElementById('project-event-item-tpl')

  // clone all other templates
  let nav_clone = document.importNode(document.getElementById('project-nav-tpl').content, true)
  let content_clone = document.importNode(document.getElementById('project-history-content-tpl').content, true)
  let event_list_clone = document.importNode(document.getElementById('project-events-list-tpl').content, true)


  // early exit if we do not have projects to display:
  if (!project.events || project.events.length < 1) {
    content_clone.getElementById('project-events-list').textContent = "No Events to display"
  }
  else {
      // fill project item template and add it to project list template
      for (item of project.events) {
          // clone tpl
          var item_clone = document.importNode(projectEventItemTpl.content, true);

          // fill with values
          item_clone.querySelector('.event-action-val').value = item.code
          item_clone.querySelector('.event-timestamp-val').value = item.timestamp

          // append to wrapper element
          event_list_clone.getElementById('project-item').appendChild(item_clone)
      }

      // project list to project list_wrapper
      content_clone.getElementById('project-list').appendChild(list_clone)
}
/*******************************                   TODO   */
// now activate the create project form
var sel = statusSelect().getElementById('status-select')
sel.name = "newProjectStatus"
form = document.getElementById('add-project-form')
form.insertBefore(sel, document.getElementById('add-project-form-submit'))
form.addEventListener('submit', addProject);

// and display everything in the page
app.textContent = ''
app.appendChild(overview_clone)
}
