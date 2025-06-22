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
    
    // New elements for edit/remove functionality
    const editRemoveBtn = document.getElementById('edit-remove-btn');
    const editRemoveModal = document.getElementById('edit-remove-modal');
    const removeCardBtn = document.getElementById('remove-card-btn');
    const cancelRemoveBtn = document.getElementById('cancel-remove-btn');

    let deck = {};
    let boxes = [[], [], [], [], []];
    let graduatedCards = [];
    let currentCardIndex = -1;
    let currentBox = 0;
    let cardConfidence = {}; // Track confidence level for each card

    const fetchDeck = async () => {
        try {
            const response = await fetch('/deck');
            deck = await response.json();
            boxes[0] = deck.cards.map((_, index) => index);
            
            // Initialize confidence levels for new cards
            deck.cards.forEach((_, index) => {
                if (!cardConfidence.hasOwnProperty(index)) {
                    cardConfidence[index] = 0; // 0 means no rating yet
                }
            });
            
            nextCard();
        } catch (error) {
            console.error("Failed to load deck:", error);
        }
    };

    const nextCard = () => {
        flashcard.classList.remove('is-flipped');
        answerControls.style.display = 'none';

        const nonEmptyBoxes = boxes
            .map((box, index) => ({ box, index }))
            .filter(item => item.box.length > 0);

        if (nonEmptyBoxes.length === 0) {
            if (graduatedCards.length === deck.cards.length) {
                cardFront.textContent = "You've completed the deck!";
                cardBack.textContent = "Congratulations!";
            } else {
                cardFront.textContent = "No cards available in boxes.";
                cardBack.textContent = "Check graduated cards or adjust logic.";
            }
            answerControls.style.display = 'none';
            updateProgressBar();
            return;
        }

        const weightedBoxSelection = [];
        nonEmptyBoxes.forEach(({ index }) => {
            const weight = (boxes.length - index) * (boxes.length - index);
            for (let i = 0; i < weight; i++) {
                weightedBoxSelection.push(index);
            }
        });

        const randomWeightedIndex = Math.floor(Math.random() * weightedBoxSelection.length);
        currentBox = weightedBoxSelection[randomWeightedIndex];

        const cardPool = boxes[currentBox];
        const randomIndexInBox = Math.floor(Math.random() * cardPool.length);
        currentCardIndex = cardPool[randomIndexInBox];

        updateCardView();
    };

    const updateCardView = () => {
        if (currentCardIndex === -1) return;
        const card = deck.cards[currentCardIndex];
        cardFront.textContent = card.front;
        cardBack.textContent = card.back;
        
        // Update confidence indicator
        const confidenceIndicator = document.getElementById('confidence-indicator');
        confidenceIndicator.className = 'confidence-indicator';
        if (cardConfidence[currentCardIndex] > 0) {
            confidenceIndicator.classList.add(`confidence-${cardConfidence[currentCardIndex]}`);
        }
        
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
        const cardIndexInBox = boxes[currentBox].indexOf(currentCardIndex);
        if (cardIndexInBox === -1) {
            console.error("Card to move not found in the current box.");
            nextCard();
            return;
        }
        
        // Store the confidence level for this card
        cardConfidence[currentCardIndex] = rating;
        
        const [cardToMove] = boxes[currentBox].splice(cardIndexInBox, 1);

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

    // Edit/Remove button functionality
    editRemoveBtn.addEventListener('click', (e) => {
        e.stopPropagation(); // Prevent card flip
        editRemoveModal.style.display = 'block';
    });

    removeCardBtn.addEventListener('click', () => {
        if (currentCardIndex !== -1) {
            // Remove card from all boxes and graduated cards
            boxes.forEach(box => {
                const index = box.indexOf(currentCardIndex);
                if (index > -1) {
                    box.splice(index, 1);
                }
            });
            
            const graduatedIndex = graduatedCards.indexOf(currentCardIndex);
            if (graduatedIndex > -1) {
                graduatedCards.splice(graduatedIndex, 1);
            }

            // Remove confidence data for this card
            delete cardConfidence[currentCardIndex];

            // Remove from deck
            deck.cards.splice(currentCardIndex, 1);

            // Update indices for remaining cards
            boxes.forEach(box => {
                for (let i = 0; i < box.length; i++) {
                    if (box[i] > currentCardIndex) {
                        box[i]--;
                    }
                }
            });

            graduatedCards.forEach((cardIndex, i) => {
                if (cardIndex > currentCardIndex) {
                    graduatedCards[i]--;
                }
            });

            // Update confidence indices for remaining cards
            const newCardConfidence = {};
            Object.keys(cardConfidence).forEach(key => {
                const oldIndex = parseInt(key);
                if (oldIndex < currentCardIndex) {
                    newCardConfidence[oldIndex] = cardConfidence[oldIndex];
                } else if (oldIndex > currentCardIndex) {
                    newCardConfidence[oldIndex - 1] = cardConfidence[oldIndex];
                }
            });
            cardConfidence = newCardConfidence;

            // Save updated deck
            saveDeck();

            // Close modal and move to next card
            editRemoveModal.style.display = 'none';
            nextCard();
        }
    });

    cancelRemoveBtn.addEventListener('click', () => {
        editRemoveModal.style.display = 'none';
    });

    // Close modals when clicking outside or on close button
    window.addEventListener('click', (e) => {
        if (e.target === questionsModal) {
            questionsModal.style.display = 'none';
        }
        if (e.target === editRemoveModal) {
            editRemoveModal.style.display = 'none';
        }
    });

    const saveDeck = async () => {
        try {
            await fetch('/deck', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(deck)
            });
        } catch (error) {
            console.error("Failed to save deck:", error);
        }
    };

    fetchDeck();
}); 