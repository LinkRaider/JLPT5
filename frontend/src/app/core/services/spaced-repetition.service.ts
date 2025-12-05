import { Injectable } from '@angular/core';

export enum ReviewQuality {
  BLACKOUT = 0,           // Complete blackout
  INCORRECT_EASY = 1,     // Incorrect, but upon seeing answer it felt easy
  INCORRECT_REMEMBERED = 2, // Incorrect, but upon seeing answer it was remembered
  CORRECT_HARD = 3,       // Correct response, but with serious difficulty
  CORRECT_HESITANT = 4,   // Correct response, but with hesitation
  PERFECT = 5             // Perfect response
}

export interface SpacedRepetitionData {
  easeFactor: number;
  interval: number;
  repetitions: number;
  nextReviewDate: Date;
}

@Injectable({
  providedIn: 'root'
})
export class SpacedRepetitionService {
  private readonly DEFAULT_EASE_FACTOR = 2.5;
  private readonly MIN_EASE_FACTOR = 1.3;

  constructor() {}

  /**
   * Calculates the next review data based on SM-2 algorithm
   * @param currentData Current spaced repetition data
   * @param quality Review quality (0-5)
   * @returns Updated spaced repetition data
   */
  calculateNextReview(currentData: SpacedRepetitionData, quality: ReviewQuality): SpacedRepetitionData {
    let easeFactor = currentData.easeFactor;
    let interval = currentData.interval;
    let repetitions = currentData.repetitions;

    // Calculate new ease factor
    // EF' = EF + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
    easeFactor = easeFactor + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02));

    // Ensure ease factor doesn't go below minimum
    if (easeFactor < this.MIN_EASE_FACTOR) {
      easeFactor = this.MIN_EASE_FACTOR;
    }

    // If quality < 3, reset repetitions and interval
    if (quality < ReviewQuality.CORRECT_HARD) {
      repetitions = 0;
      interval = 1;
    } else {
      repetitions++;

      // Calculate new interval based on repetitions
      if (repetitions === 1) {
        interval = 1;
      } else if (repetitions === 2) {
        interval = 6;
      } else {
        interval = Math.round(interval * easeFactor);
      }
    }

    // Calculate next review date
    const nextReviewDate = new Date();
    nextReviewDate.setDate(nextReviewDate.getDate() + interval);

    return {
      easeFactor,
      interval,
      repetitions,
      nextReviewDate
    };
  }

  /**
   * Initialize spaced repetition data for a new item
   */
  initializeData(): SpacedRepetitionData {
    return {
      easeFactor: this.DEFAULT_EASE_FACTOR,
      interval: 1,
      repetitions: 0,
      nextReviewDate: new Date()
    };
  }

  /**
   * Get a human-readable description of the review quality
   */
  getQualityDescription(quality: ReviewQuality): string {
    switch (quality) {
      case ReviewQuality.BLACKOUT:
        return 'Complete blackout';
      case ReviewQuality.INCORRECT_EASY:
        return 'Incorrect, but felt easy';
      case ReviewQuality.INCORRECT_REMEMBERED:
        return 'Incorrect, but remembered';
      case ReviewQuality.CORRECT_HARD:
        return 'Correct with difficulty';
      case ReviewQuality.CORRECT_HESITANT:
        return 'Correct with hesitation';
      case ReviewQuality.PERFECT:
        return 'Perfect!';
      default:
        return 'Unknown';
    }
  }

  /**
   * Get button labels for quality ratings
   */
  getQualityButtonLabels(): Array<{ quality: ReviewQuality; label: string; color: string }> {
    return [
      { quality: ReviewQuality.BLACKOUT, label: 'Again', color: 'danger' },
      { quality: ReviewQuality.CORRECT_HARD, label: 'Hard', color: 'warning' },
      { quality: ReviewQuality.CORRECT_HESITANT, label: 'Good', color: 'primary' },
      { quality: ReviewQuality.PERFECT, label: 'Easy', color: 'success' }
    ];
  }
}
