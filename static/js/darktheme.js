let darkMode = localStorage.getItem('darkMode'); 

const darkModeToggle = document.querySelector('#dark-mode-toggle');

function enableDarkMode () {
  document.getElementById("maincss").href = "/static/css/main-dark.css";
  localStorage.setItem('darkMode', 'enabled');
}

function disableDarkMode () {
  document.getElementById("maincss").href = "/static/css/main.css";
  localStorage.setItem('darkMode', null);
}
 
if (darkMode === 'enabled') {
  enableDarkMode();
}

darkModeToggle.addEventListener('click', () => {
  darkMode = localStorage.getItem('darkMode'); 
  
  if (darkMode !== 'enabled') {
    enableDarkMode();
  } else {  
    disableDarkMode(); 
  }
});
