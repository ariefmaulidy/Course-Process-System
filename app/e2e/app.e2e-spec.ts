import { ProjekmpplPage } from './app.po';

describe('projekmppl App', () => {
  let page: ProjekmpplPage;

  beforeEach(() => {
    page = new ProjekmpplPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
