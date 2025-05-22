import {Card} from './CardComponent/card.js'

document.addEventListener('DOMContentLoaded',()=>{
    const root=document.getElementById("message_layout");
    root.style.height='100px'
    root.style.display='flex';
    root.style.flexDirection='column';
    root.style.background='var(--color-white)';
    root.style.padding='var(--card-padding)';
    root.style.borderRadius='var(--card-border-radius)';
    root.style.fontSize='1.4rem';
    root.style.height='max-content';  
    root.appendChild(Card());
})