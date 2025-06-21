document.addEventListener('DOMContentLoaded', () => {
    const flashcard = document.querySelector('.flashcard');
    const cardFront = document.querySelector('.card-front');
    const cardBack = document.querySelector('.card-back');
    const incorrectBtn = document.getElementById('incorrect-btn');
    const correctBtn = document.getElementById('correct-btn');
    const progressIndicator = document.getElementById('progress-indicator');

    let deck = {};
    let boxes = [[], [], []]; // 3 boxes for the Leitner system
    let currentCardIndex = -1;
    let currentBox = 0;

    const fetchDeck = async () => {
        const response = await fetch('/deck');
        deck = await response.json();
        boxes[0] = deck.cards.map((_, index) => index);
        nextCard();
    };

    const nextCard = () => {
        flashcard.classList.remove('is-flipped');
        // Find the next card to show
        for (let i = 0; i < boxes.length; i++) {
            if (boxes[i].length > 0) {
                currentBox = i;
                currentCardIndex = boxes[i][0];
                updateCardView();
                return;
            }
        }
        // No cards left
        cardFront.textContent = "You've completed the deck!";
        cardBack.textContent = "Congratulations!";
        incorrectBtn.style.display = 'none';
        correctBtn.style.display = 'none';
    };

    const updateCardView = () => {
        const card = deck.cards[currentCardIndex];
        cardFront.textContent = card.front;
        cardBack.textContent = card.back;
        progressIndicator.textContent = `Box 1: ${boxes[0].length}, Box 2: ${boxes[1].length}, Box 3: ${boxes[2].length}`;
    };

    flashcard.addEventListener('click', () => {
        flashcard.classList.toggle('is-flipped');
    });

    correctBtn.addEventListener('click', () => {
        boxes[currentBox].shift();
        if (currentBox < boxes.length - 1) {
            boxes[currentBox + 1].push(currentCardIndex);
        }
        nextCard();
    });

    incorrectBtn.addEventListener('click', () => {
        boxes[currentBox].shift();
        boxes[0].push(currentCardIndex); // Move back to the first box
        nextCard();
    });

    fetchDeck();
}); 