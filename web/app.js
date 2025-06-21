document.addEventListener('DOMContentLoaded', () => {
    const deckTitle = document.getElementById('deck-title');
    const cardList = document.getElementById('card-list');
    const newCardForm = document.getElementById('new-card-form');
    const newFront = document.getElementById('new-front');
    const newBack = document.getElementById('new-back');
    const saveBtn = document.getElementById('save-btn');

    let deck = {};

    const fetchDeck = async () => {
        try {
            const response = await fetch('/deck');
            deck = await response.json();
            renderDeck();
        } catch (error) {
            console.error('Failed to load deck:', error);
        }
    };

    const renderDeck = () => {
        deckTitle.textContent = `Editing: ${deck.name}`;
        cardList.innerHTML = '';
        deck.cards.forEach((card, index) => {
            const cardEl = document.createElement('div');
            cardEl.className = 'card';
            cardEl.innerHTML = `
                <textarea class="front" data-index="${index}">${card.front}</textarea>
                <textarea class="back" data-index="${index}">${card.back}</textarea>
                <div class="card-actions">
                    <button class="btn-danger delete-btn" data-index="${index}">Delete</button>
                </div>
            `;
            cardList.appendChild(cardEl);
        });
    };

    newCardForm.addEventListener('submit', (e) => {
        e.preventDefault();
        deck.cards.push({ front: newFront.value, back: newBack.value });
        newFront.value = '';
        newBack.value = '';
        renderDeck();
    });

    cardList.addEventListener('input', (e) => {
        const index = e.target.dataset.index;
        if (e.target.classList.contains('front')) {
            deck.cards[index].front = e.target.value;
        } else if (e.target.classList.contains('back')) {
            deck.cards[index].back = e.target.value;
        }
    });

    cardList.addEventListener('click', (e) => {
        if (e.target.classList.contains('delete-btn')) {
            const index = e.target.dataset.index;
            deck.cards.splice(index, 1);
            renderDeck();
        }
    });

    saveBtn.addEventListener('click', async () => {
        try {
            await fetch('/deck', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(deck),
            });
            alert('Deck saved successfully!');
        } catch (error) {
            console.error('Failed to save deck:', error);
            alert('Error saving deck.');
        }
    });

    fetchDeck();
}); 