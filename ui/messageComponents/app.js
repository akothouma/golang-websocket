import Card from './CardComponent/card.js'
document.addEventListener('DOMContentLoaded',()=>{
    const root=document.getElementById("chat_layout");
    root.id="chat_layout"
    root.appendChild(Card());
})