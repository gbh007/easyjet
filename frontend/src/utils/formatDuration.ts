/**
 * Форматирует длительность в миллисекундах в человекочитаемый формат
 * @param ms - Длительность в миллисекундах
 * @returns Отформатированная строка (например, "1м 30с", "45с", "2ч 15м 10с")
 */
export function formatDuration(ms: number): string {
  if (ms < 1000) {
    return `${ms}мс`;
  }

  const seconds = Math.floor(ms / 1000);
  if (seconds < 60) {
    return `${seconds}с`;
  }

  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  if (minutes < 60) {
    return `${minutes}м ${remainingSeconds}с`;
  }

  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  return `${hours}ч ${remainingMinutes}м ${remainingSeconds}с`;
}
