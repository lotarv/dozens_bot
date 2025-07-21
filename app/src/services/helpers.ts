export function formatDateTime(isoDate: string): string {
    const date = new Date(isoDate);
    
    
    // Получаем день, месяц и время
    const day = date.getUTCDate();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const year= date.getFullYear();
    
    return `${day}.${month}.${year}`;
}