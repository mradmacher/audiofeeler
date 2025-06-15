describe("template spec", () => {
  it("passes", () => {
    cy.visit("http://localhost:3000")
    cy.contains("Test").click()

    cy.title().should("eq", "Audiofeeler")
    cy.get("h1").should("have.text", "Test")
  })
})
