// Elements
const landing = document.getElementById("landing");
const authPage = document.getElementById("authPage");
const dashboard = document.getElementById("dashboard");
const guestBtn = document.getElementById("guestBtn");
const loginBtn = document.getElementById("loginBtn");
const signupBtn = document.getElementById("signupBtn");
const backBtn = document.getElementById("backBtn");
const dashBackBtn = document.getElementById("dashBackBtn");
const toggleAuth = document.getElementById("toggleAuth");
const authForm = document.getElementById("authForm");
const authTitle = document.getElementById("authTitle");
const usernameInput = document.getElementById("username");
const passwordInput = document.getElementById("password");

const transactionForm = document.getElementById("transactionForm");
const descInput = document.getElementById("desc");
const amountInput = document.getElementById("amount");
const categoryInput = document.getElementById("category");
const list = document.getElementById("list");

const balanceEl = document.getElementById("balance");
const incomeEl = document.getElementById("income");
const expensesEl = document.getElementById("expenses");
const savingsEl = document.getElementById("savings");
const emergencyEl = document.getElementById("emergency");
const logoutBtn = document.getElementById("logoutBtn");

let currentAuth = "login";
let token = localStorage.getItem("token");
let username = localStorage.getItem("username");
let transactions = JSON.parse(localStorage.getItem("transactions") || "[]");

// Show Page
function showPage(page){
  landing.classList.add("hidden");
  authPage.classList.add("hidden");
  dashboard.classList.add("hidden");
  page.classList.remove("hidden");
}

// Landing buttons
guestBtn.onclick = () => { showPage(dashboard); renderTransactions(); }
loginBtn.onclick = () => { showPage(authPage); authTitle.innerText="Log In"; currentAuth="login"; }
signupBtn.onclick = () => { showPage(authPage); authTitle.innerText="Sign Up"; currentAuth="signup"; }

// Back buttons
backBtn.onclick = () => showPage(landing);
dashBackBtn.onclick = () => showPage(landing);

// Toggle login/signup
toggleAuth.onclick = () => {
  if(currentAuth==="login"){ currentAuth="signup"; authTitle.innerText="Sign Up"; toggleAuth.innerText="Already have an account? Log In"; }
  else{ currentAuth="login"; authTitle.innerText="Log In"; toggleAuth.innerText="Don't have an account? Sign Up"; }
}

// Auth submit
authForm.onsubmit = async (e)=>{
  e.preventDefault();
  const user={username: usernameInput.value,password: passwordInput.value};
  const url = currentAuth==="login"?"/login":"/signup";
  try{
    const res = await fetch(url,{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify(user)});
    if(!res.ok){ alert("Error: "+res.statusText); return; }
    if(currentAuth==="login"){ const data=await res.json(); token=data.token; localStorage.setItem("token",token); localStorage.setItem("username",user.username); username=user.username; fetchTransactions(); }
    showPage(dashboard);
  }catch(err){ console.log(err); alert("Server error"); }
};

// Logout
logoutBtn.onclick = ()=>{
  localStorage.removeItem("token");
  localStorage.removeItem("username");
  showPage(landing);
};

// Transactions
transactionForm.onsubmit = async (e)=>{
  e.preventDefault();
  const t={desc:descInput.value,amount:parseFloat(amountInput.value),category:categoryInput.value};
  if(token){
    await fetch("/transactions",{method:"POST",headers:{"Content-Type":"application/json","Authorization":token},body:JSON.stringify(t)});
    await fetchTransactions();
  } else {
    transactions.push(t);
    localStorage.setItem("transactions",JSON.stringify(transactions));
    renderTransactions();
  }
  descInput.value=""; amountInput.value=""; categoryInput.value="Income";
};

// Render transactions and balances
function renderTransactions(){
  const tx = token? window.txList : transactions;
  let income=0, expenses=0, savings=0, emergency=0, balance=0;
  list.innerHTML="";
  tx.forEach(t=>{
    const li=document.createElement("li");
    li.innerText = `${t.desc} (${t.category}): Ksh ${t.amount}`;
    list.appendChild(li);
    switch(t.category){
      case "Income": income+=t.amount; break;
      case "Expense": expenses+=t.amount; break;
      case "Saving": savings+=t.amount; break;
      case "Emergency": emergency+=t.amount; break;
    }
  });
  balance = income - expenses - savings - emergency;
  balanceEl.innerText=`Ksh ${balance}`;
  incomeEl.innerText=`Ksh ${income}`;
  expensesEl.innerText=`Ksh ${expenses}`;
  savingsEl.innerText=`Ksh ${savings}`;
  emergencyEl.innerText=`Ksh ${emergency}`;
}

// Fetch backend transactions
async function fetchTransactions(){
  try{
    const res=await fetch("/transactions",{headers:{"Authorization":token}});
    const data=await res.json();
    window.txList=data;
    renderTransactions();
  }catch(err){ console.log(err); alert("Cannot fetch transactions"); }
}

// Edit buttons
document.querySelectorAll(".edit-btn").forEach(btn=>{
  btn.onclick = ()=>{
    const type = btn.dataset.type;
    let current = 0;
    switch(type){
      case "Income": current = parseFloat(incomeEl.innerText.replace("Ksh ","")); break;
      case "Saving": current = parseFloat(savingsEl.innerText.replace("Ksh ","")); break;
      case "Emergency": current = parseFloat(emergencyEl.innerText.replace("Ksh ","")); break;
      case "Expense": current = parseFloat(expensesEl.innerText.replace("Ksh ","")); break; // Added Expense
    }
    const newAmount = parseFloat(prompt(`Enter new amount for ${type}:`, current));
    if(isNaN(newAmount)) return alert("Invalid amount");

    if(token){
      const tx = window.txList.find(t=>t.category===type);
      if(tx){
        tx.amount = newAmount;
        fetch("/transactions",{method:"POST",headers:{"Content-Type":"application/json","Authorization":token},body:JSON.stringify(tx)})
        .then(()=> fetchTransactions());
      } else {
        const t = {desc: type+" adjustment", amount:newAmount, category:type};
        fetch("/transactions",{method:"POST",headers:{"Content-Type":"application/json","Authorization":token},body:JSON.stringify(t)})
        .then(()=> fetchTransactions());
      }
    } else {
      const index = transactions.findIndex(t=>t.category===type);
      if(index!==-1){ transactions[index].amount=newAmount; }
      else{ transactions.push({desc:type+" adjustment", amount:newAmount, category:type}); }
      localStorage.setItem("transactions",JSON.stringify(transactions));
      renderTransactions();
    }
  };
});

// Init page
if(token){ showPage(dashboard); fetchTransactions(); } else { showPage(landing); renderTransactions(); }