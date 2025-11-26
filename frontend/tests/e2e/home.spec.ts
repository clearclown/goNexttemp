import { expect, test } from "@playwright/test";

test.describe("Home Page", () => {
  test("should display the main heading", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("h1")).toContainText("goNexttemp");
  });

  test("should have login and register links", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator('a[href="/login"]')).toBeVisible();
    await expect(page.locator('a[href="/register"]')).toBeVisible();
  });

  test("should navigate to login page", async ({ page }) => {
    await page.goto("/");
    await page.click('a[href="/login"]');
    await expect(page).toHaveURL("/login");
    await expect(page.locator("h1")).toContainText("ログイン");
  });

  test("should navigate to register page", async ({ page }) => {
    await page.goto("/");
    await page.click('a[href="/register"]');
    await expect(page).toHaveURL("/register");
    await expect(page.locator("h1")).toContainText("新規登録");
  });
});
