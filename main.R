# which(...); any(...); all(...); unique(...); duplicated(...)
# is.element(x, y); x %in% y; substr(x, start, stop)
# length(x); sum(x); min(x); max(x); abs(x); sqrt(x)
# mean(x); median(x); quantile(x, p); sd(x); var(x)
# table(x); table(x, y); plot(x, y); plot( y~x, data ); boxplot(x); boxplot( y~x, data )
# barplot( table(x), beside=F, legend=F ); hist( x, probability=F, breaks )
# rep( x, times=1, length.out=NA, each=1 ); seq( from, to, by )
# sample( x, size, replace=F, prob ); replicate( n, expr )
# dbinom(k, n, p) dgeom(k, p) dnbinom(k, r, p) dpois(k, lambda) dhyper(k, M, N-M, n)
# pbinom(k, n, p) pgeom(k, p) pnbinom(k, r, p) ppois(k, lambda) phyper(k, M, N-M, n)
# punif(q, a, b) pexp(q, lambda) pnorm(q, mu, sigma) pt(q, df) pchisq(q, df)
# qunif(p, a, b) qexp(p, lambda) qnorm(p, mu, sigma) qt(p, df) qchisq(p, df)
# t.test( x, mu, alternative=c("two.sided", "less", "greater"), conf.level=0.95 )
# t.test( x, y, alternative=c("two.sided", "less", "greater"), paired=F )
# t.test(...)$p.value; t.test(...)$conf.int
# prop.test( x, n, p, alternative=c("two.sided", "less", "greater"), correct=T )
# chisq.test(x, p)
# chisq.test(x)

x <- c(28, 36, 36, 30, 27, 23)
probs <- rep(1/6, 6)

n <- sum(x)
k <- length(probs)

(x - n*probs)^2

# chi2.obs <- sum((x - n*probs)^2 / (n*probs))
# chi2.obs
# p.value <- 1-pchisq(chi2.obs, df=k-1)
# p.value
#
# chisq.test(x, p=probs)$p.value

