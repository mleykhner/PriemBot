export default async ({ page }) => {
    await page.goto('https://mai.ru/press/news/', { waitUntil: 'networkidle2' });
    await page.waitForSelector('article');

    return await page.$$eval('article', articles => articles
        .map(a => {
            const id = a.getAttribute('id')?.split('_').pop() || null;
            const img = a.querySelector('img.card-img-top');
            const title = a.querySelector('h5');
            const subtitle = a.querySelector('p.small.text-muted.mb-0');
            return {
                id: id,
                imageUrl: img ? img.src : null,
                title: title ? title.textContent.trim() : null,
                subtitle: subtitle ? subtitle.textContent.trim() : null
            };
        })
        .filter(obj =>
            obj.id !== null &&
            obj.imageUrl !== null &&
            obj.title !== null &&
            obj.subtitle !== null
        )
    );
}