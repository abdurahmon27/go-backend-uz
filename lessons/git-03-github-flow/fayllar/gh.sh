# Haqiqiy GitHub'da (bu darsda ishga tushirmaymiz)
gh pr create --base main --title "chegirma: foiz bilan narx" \
  --body-file pr.md --reviewer aziz     # PR #2 ochiladi
gh pr checkout 2                        # reviewer PR'ni oladi
gh pr review 2 --request-changes --body-file izoh.md
gh pr review 2 --approve --body "Rahmat, tayyor!"
gh pr merge 2 --merge --delete-branch   # merge commit
git fetch origin pull/2/head:pr-2       # PR ham ref
