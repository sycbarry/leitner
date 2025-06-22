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
    let focusCards = []; // Current 5-card focus set
    let focusCardIndex = 0; // Current position within focus set
    let cardsStudiedInFocus = 0; // Track how many cards studied in current focus set
    let lastResurfaceTime = 0; // Track when we last resurfaced a well-performed card

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
            
            initializeFocusSet();
            nextCard();
        } catch (error) {
            console.error("Failed to load deck:", error);
        }
    };

    const initializeFocusSet = () => {
        // Get all available cards (not graduated)
        const availableCards = [];
        boxes.forEach((box, boxIndex) => {
            box.forEach(cardIndex => {
                availableCards.push({ cardIndex, boxIndex });
            });
        });

        if (availableCards.length === 0) {
            focusCards = [];
            return;
        }

        // Prioritize lower boxes (harder cards) for focus set
        availableCards.sort((a, b) => {
            // First prioritize by box (lower = harder)
            if (a.boxIndex !== b.boxIndex) {
                return a.boxIndex - b.boxIndex;
            }
            // Then by confidence (lower = needs more work)
            const confidenceA = cardConfidence[a.cardIndex] || 0;
            const confidenceB = cardConfidence[b.cardIndex] || 0;
            return confidenceA - confidenceB;
        });

        // Take up to 5 cards for focus set
        focusCards = availableCards.slice(0, 5).map(item => item.cardIndex);
        focusCardIndex = 0;
        cardsStudiedInFocus = 0;
    };

    const shouldResurfaceWellPerformedCard = () => {
        // Resurface a well-performed card every 3-5 cards studied
        const resurfaceInterval = Math.random() < 0.3 ? 3 : 5;
        return cardsStudiedInFocus >= resurfaceInterval;
    };

    const getWellPerformedCard = () => {
        // Find cards with confidence 4 or 5 that are not in current focus set
        const wellPerformedCards = [];
        boxes.forEach((box, boxIndex) => {
            box.forEach(cardIndex => {
                const confidence = cardConfidence[cardIndex] || 0;
                if (confidence >= 4 && !focusCards.includes(cardIndex)) {
                    wellPerformedCards.push(cardIndex);
                }
            });
        });

        // Also include graduated cards as well-performed
        graduatedCards.forEach(cardIndex => {
            if (!focusCards.includes(cardIndex)) {
                wellPerformedCards.push(cardIndex);
            }
        });

        if (wellPerformedCards.length === 0) {
            return null;
        }

        // Return a random well-performed card
        return wellPerformedCards[Math.floor(Math.random() * wellPerformedCards.length)];
    };

    const nextCard = () => {
        flashcard.classList.remove('is-flipped');
        answerControls.style.display = 'none';

        // Check if we should resurface a well-performed card
        if (shouldResurfaceWellPerformedCard()) {
            const resurfaceCard = getWellPerformedCard();
            if (resurfaceCard !== null) {
                currentCardIndex = resurfaceCard;
                currentBox = getCardBox(resurfaceCard);
                cardsStudiedInFocus = 0; // Reset counter
                updateCardView();
                return;
            }
        }

        // Check if we need to refresh the focus set
        if (focusCards.length === 0 || focusCardIndex >= focusCards.length) {
            initializeFocusSet();
            if (focusCards.length === 0) {
                // No more cards to study
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
        }

        // Get next card from focus set
        currentCardIndex = focusCards[focusCardIndex];
        currentBox = getCardBox(currentCardIndex);
        focusCardIndex++;
        cardsStudiedInFocus++;
        updateCardView();
    };

    const getCardBox = (cardIndex) => {
        for (let i = 0; i < boxes.length; i++) {
            if (boxes[i].includes(cardIndex)) {
                return i;
            }
        }
        return -1; // Card not found in any box
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
            // Refresh focus set since we graduated a card
            initializeFocusSet();
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
        
        // If the card moved significantly (up or down), refresh focus set
        if (Math.abs(newBox - currentBox) >= 2) {
            initializeFocusSet();
        }
        
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

            // Refresh focus set since we removed a card
            initializeFocusSet();

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