document.addEventListener('DOMContentLoaded', () => {
    const flashcard = document.querySelector('.flashcard');
    const cardFront = document.querySelector('.card-front');
    const cardBack = document.querySelector('.card-back');
    const answerControls = document.getElementById('answer-controls');
    const progressBar = document.getElementById('progress-bar');
    const viewAllBtn = document.getElementById('view-all-btn');
    const questionsModal = document.getElementById('questions-modal');
    const closeModalBtn = document.querySelector('.close-btn');
    const allQuestionsList = document.getElementById('all-questions-list');

    let deck = {};
    let boxes = [[], [], [], [], []];
    let graduatedCards = [];
    let currentCardIndex = -1;
    let currentBox = 0;

    const fetchDeck = async () => {
        try {
            const response = await fetch('/deck');
            deck = await response.json();
            boxes[0] = deck.cards.map((_, index) => index);
            nextCard();
        } catch (error) {
            console.error("Failed to load deck:", error);
        }
    };

    const nextCard = () => {
        flashcard.classList.remove('is-flipped');
        answerControls.style.display = 'none';

        for (let i = 0; i < boxes.length; i++) {
            if (boxes[i].length > 0) {
                currentBox = i;
                currentCardIndex = boxes[i][0];
                updateCardView();
                return;
            }
        }

        cardFront.textContent = "You've completed the deck!";
        cardBack.textContent = "Congratulations!";
        answerControls.style.display = 'none';
        updateProgressBar(); // Final update to 100%
    };

    const updateCardView = () => {
        if (currentCardIndex === -1) return;
        const card = deck.cards[currentCardIndex];
        cardFront.textContent = card.front;
        cardBack.textContent = card.back;
        updateProgressBar();
    };

    const updateProgressBar = () => {
        const totalCards = deck.cards.length;
        if (totalCards === 0) {
            progressBar.style.height = '0%';
            return;
        }

        const maxScore = totalCards * (boxes.length); // +1 point for graduating
        let currentScore = 0;
        boxes.forEach((box, boxIndex) => {
            currentScore += box.length * boxIndex;
        });
        currentScore += graduatedCards.length * (boxes.length);

        const progressPercentage = (currentScore / maxScore) * 100;
        progressBar.style.height = `${progressPercentage}%`;
    };

    flashcard.addEventListener('click', () => {
        if (currentCardIndex === -1) return;
        flashcard.classList.toggle('is-flipped');
        answerControls.style.display = flashcard.classList.contains('is-flipped') ? 'block' : 'none';
    });

    answerControls.addEventListener('click', (e) => {
        if (e.target.matches('.btn-rating')) {
            const rating = parseInt(e.target.dataset.rating, 10);
            moveCard(rating);
        }
    });

    const moveCard = (rating) => {
        const cardToMove = boxes[currentBox].shift();

        if (currentBox === boxes.length - 1 && rating >= 4) {
            graduatedCards.push(cardToMove);
            nextCard();
            return;
        }

        let newBox;
        switch (rating) {
            case 1: newBox = 0; break;
            case 2: newBox = 0; break;
            case 3: newBox = Math.min(currentBox + 1, boxes.length - 1); break;
            case 4: newBox = Math.min(currentBox + 2, boxes.length - 1); break;
            case 5: newBox = Math.min(currentBox + 3, boxes.length - 1); break;
            default: newBox = currentBox;
        }
        
        // Ensure "Hard" doesn't demote a card
        if (rating === 3 && newBox < currentBox) {
            newBox = currentBox;
        }

        boxes[newBox].push(cardToMove);
        nextCard();
    };

    const populateQuestionsList = () => {
        allQuestionsList.innerHTML = '';
        deck.cards.forEach(card => {
            const li = document.createElement('li');
            li.textContent = card.front;
            allQuestionsList.appendChild(li);
        });
    };

    viewAllBtn.addEventListener('click', () => {
        populateQuestionsList();
        questionsModal.style.display = 'block';
    });

    closeModalBtn.addEventListener('click', () => {
        questionsModal.style.display = 'none';
    });

    window.addEventListener('click', (e) => {
        if (e.target === questionsModal) {
            questionsModal.style.display = 'none';
        }
    });

    fetchDeck();
}); 