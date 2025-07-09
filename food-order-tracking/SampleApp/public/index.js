// Add new menu item
async function addMenu(e) {
    e.preventDefault();
    const menuID = document.getElementById("menuID").value;
    const foodName = document.getElementById("foodName").value;
    const price = parseInt(document.getElementById("price").value);
  
    const res = await fetch("/api/menu", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ menuID, foodName, price })
    });
  
    try {
      const data = await res.json();
      alert("Menu added successfully!");
      document.getElementById("menuForm").reset();
      getAllMenus();
    } catch (e) {
      alert("Invalid response from server");
    }
  }
  
  // Fetch and display all menu items
  async function getAllMenus() {
    const res = await fetch("/api/menu/all");
    const data = await res.json();
    displayMenus(data.data || []);
  }
  
  // Fetch and display menu by ID
  async function getMenuByID() {
    const id = document.getElementById("menuIdInput").value;
    if (!id) return alert("Enter Menu ID");
    const res = await fetch(`/api/menu/${id}`);
    const data = await res.json();
    if (data.data) displayMenus([data.data]);
    else alert("Menu not found");
  }
  
  // Search menu by food name
  async function searchMenu() {
    const name = document.getElementById("foodNameInput").value;
    if (!name) return alert("Enter Food Name");
    const res = await fetch(`/api/menu/search?name=${name}`);
    const data = await res.json();
    displayMenus(data.data || []);
  }
  
  // Display menus in HTML
  function displayMenus(menus) {
    const tableBody = document.getElementById("menuTableBody");
    tableBody.innerHTML = "";
    menus.forEach((menu) => {
      const row = `<tr>
          <td>${menu.menuID}</td>
          <td>${menu.foodName}</td>
          <td>${menu.price}</td>
        </tr>`;
      tableBody.innerHTML += row;
    });
  }
  
  // Event Listeners
  document.getElementById("menuForm").addEventListener("submit", addMenu);
  document.getElementById("btnAllMenus").addEventListener("click", getAllMenus);
  document.getElementById("btnGetMenu").addEventListener("click", getMenuByID);
  document.getElementById("btnSearch").addEventListener("click", searchMenu);
  