package parser

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Permify/permify/pkg/dsl/ast"
)

// TestParser -
func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "parser-suite")
}

var _ = Describe("parser", func() {
	Context("Statement", func() {
		It("Case 1 - Repository with parent and owner relations and read action", func() {
			pr := NewParser(`
			entity repository {
		
			relation parent @organization
			relation owner  @user
		
			action read = owner and (parent.admin and not parent.member)
		
			}`)

			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())
			st := schema.Statements[0].(*ast.EntityStatement)

			Expect(st.Name.Literal).Should(Equal("repository"))

			r1 := st.RelationStatements[0].(*ast.RelationStatement)
			Expect(r1.Name.Literal).Should(Equal("parent"))

			for _, a := range r1.RelationTypes {
				Expect(a.Type.Literal).Should(Equal("organization"))
			}

			r2 := st.RelationStatements[1].(*ast.RelationStatement)
			Expect(r2.Name.Literal).Should(Equal("owner"))

			for _, a := range r2.RelationTypes {
				Expect(a.Type.Literal).Should(Equal("user"))
			}

			a1 := st.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(a1.Name.Literal).Should(Equal("read"))

			es := a1.ExpressionStatement.(*ast.ExpressionStatement)

			Expect(es.Expression.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("owner"))
			Expect(es.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("parent.admin"))
			Expect(es.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("not parent.member"))
		})

		It("Case 2 - Repository with parent and owner relations and read action", func() {
			pr := NewParser(`
			entity repository {
				relation parent   @organization
				relation owner  @user
		
				action read = (owner and parent.admin) and parent.member
			}`)

			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())
			st := schema.Statements[0].(*ast.EntityStatement)

			Expect(st.Name.Literal).Should(Equal("repository"))

			r1 := st.RelationStatements[0].(*ast.RelationStatement)
			Expect(r1.Name.Literal).Should(Equal("parent"))

			for _, a := range r1.RelationTypes {
				Expect(a.Type.Literal).Should(Equal("organization"))
			}

			r2 := st.RelationStatements[1].(*ast.RelationStatement)
			Expect(r2.Name.Literal).Should(Equal("owner"))

			for _, a := range r2.RelationTypes {
				Expect(a.Type.Literal).Should(Equal("user"))
			}

			a1 := st.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(a1.Name.Literal).Should(Equal("read"))

			es := a1.ExpressionStatement.(*ast.ExpressionStatement)

			Expect(es.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("owner"))
			Expect(es.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("parent.admin"))
			Expect(es.Expression.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("parent.member"))
		})

		It("Case 3 - Organization with owner relation and delete action", func() {
			pr := NewParser(`
			entity organization {
				relation owner @user
				action delete = owner
			}
			`)
			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())
			st := schema.Statements[0].(*ast.EntityStatement)

			Expect(st.Name.Literal).Should(Equal("organization"))

			r1 := st.RelationStatements[0].(*ast.RelationStatement)
			Expect(r1.Name.Literal).Should(Equal("owner"))

			for _, a := range r1.RelationTypes {
				Expect(a.Type.Literal).Should(Equal("user"))
			}

			a1 := st.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(a1.Name.Literal).Should(Equal("delete"))

			es := a1.ExpressionStatement.(*ast.ExpressionStatement)

			Expect(es.Expression.(*ast.Identifier).String()).Should(Equal("owner"))
		})

		It("Case 4: Organization with owner relation and delete action", func() {
			pr := NewParser("entity organization {\n\nrelation owner @user\n\naction delete = not owner\n\n\n}\n\n")
			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())

			st := schema.Statements[0].(*ast.EntityStatement)

			Expect(st.Name.Literal).Should(Equal("organization"))

			r1 := st.RelationStatements[0].(*ast.RelationStatement)
			Expect(r1.Name.Literal).Should(Equal("owner"))

			for _, a := range r1.RelationTypes {
				Expect(a.Type.Literal).Should(Equal("user"))
			}

			a1 := st.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(a1.Name.Literal).Should(Equal("delete"))

			es := a1.ExpressionStatement.(*ast.ExpressionStatement)

			Expect(es.Expression.(*ast.Identifier).String()).Should(Equal("not owner"))
		})

		It("Case 5 - Repository view and read actions with ownership and parent organization", func() {
			pr := NewParser(`
			entity repository {
		
				relation parent  @organization
				relation owner  @user @organization#member
		
				action view = owner
				action read = view and (parent.admin and parent.member)
			}
			`)

			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())

			st := schema.Statements[0].(*ast.EntityStatement)

			Expect(st.Name.Literal).Should(Equal("repository"))

			r1 := st.RelationStatements[0].(*ast.RelationStatement)
			Expect(r1.Name.Literal).Should(Equal("parent"))

			for _, a := range r1.RelationTypes {
				Expect(a.Type.Literal).Should(Equal("organization"))
			}

			r2 := st.RelationStatements[1].(*ast.RelationStatement)
			Expect(r2.Name.Literal).Should(Equal("owner"))

			Expect(r2.RelationTypes[0].Type.Literal).Should(Equal("user"))
			Expect(r2.RelationTypes[1].Type.Literal).Should(Equal("organization"))
			Expect(r2.RelationTypes[1].Relation.Literal).Should(Equal("member"))

			a1 := st.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(a1.Name.Literal).Should(Equal("view"))

			es1 := a1.ExpressionStatement.(*ast.ExpressionStatement)
			Expect(es1.Expression.(*ast.Identifier).String()).Should(Equal("owner"))

			a2 := st.PermissionStatements[1].(*ast.PermissionStatement)
			Expect(a2.Name.Literal).Should(Equal("read"))

			es2 := a2.ExpressionStatement.(*ast.ExpressionStatement)

			Expect(es2.Expression.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("view"))
			Expect(es2.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("parent.admin"))
			Expect(es2.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("parent.member"))
		})

		It("Case 6 - Complex organization and repository permissions", func() {
			pr := NewParser(`
			entity user {}

			entity organization {
    			// relations
				relation admin @user
    			relation member @user

				// actions
    			action create_repository = (admin or member)
			}

			entity repository {
    			// relations
    			relation owner @user @organization#member
    			relation parent @organization
    
    			// actions
    			permission read = (owner and (parent.admin and not parent.member)) or owner
    
    			// parent.create_repository means user should be
    			// organization admin or organization member
    			permission delete = (owner or (parent.create_repository))
			}
			`)

			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())

			// USER
			userSt := schema.Statements[0].(*ast.EntityStatement)
			Expect(userSt.Name.Literal).Should(Equal("user"))

			// ORGANIZATION
			organizationSt := schema.Statements[1].(*ast.EntityStatement)

			Expect(organizationSt.Name.Literal).Should(Equal("organization"))

			or1 := organizationSt.RelationStatements[0].(*ast.RelationStatement)
			Expect(or1.Name.Literal).Should(Equal("admin"))

			Expect(or1.RelationTypes[0].Type.Literal).Should(Equal("user"))

			or2 := organizationSt.RelationStatements[1].(*ast.RelationStatement)
			Expect(or2.Name.Literal).Should(Equal("member"))

			Expect(or2.RelationTypes[0].Type.Literal).Should(Equal("user"))

			oa1 := organizationSt.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(oa1.Name.Literal).Should(Equal("create_repository"))

			oes1 := oa1.ExpressionStatement.(*ast.ExpressionStatement)

			Expect(oes1.Expression.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("admin"))
			Expect(oes1.Expression.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("member"))

			// REPOSITORY

			repositorySt := schema.Statements[2].(*ast.EntityStatement)

			Expect(repositorySt.Name.Literal).Should(Equal("repository"))

			rr1 := repositorySt.RelationStatements[0].(*ast.RelationStatement)
			Expect(rr1.Name.Literal).Should(Equal("owner"))

			Expect(rr1.RelationTypes[0].Type.Literal).Should(Equal("user"))
			Expect(rr1.RelationTypes[1].Type.Literal).Should(Equal("organization"))
			Expect(rr1.RelationTypes[1].Relation.Literal).Should(Equal("member"))

			rr2 := repositorySt.RelationStatements[1].(*ast.RelationStatement)
			Expect(rr2.Name.Literal).Should(Equal("parent"))

			Expect(rr2.RelationTypes[0].Type.Literal).Should(Equal("organization"))

			ra1 := repositorySt.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(ra1.Name.Literal).Should(Equal("read"))

			res1 := ra1.ExpressionStatement.(*ast.ExpressionStatement)

			Expect(res1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("owner"))
			Expect(res1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("parent.admin"))
			Expect(res1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("not parent.member"))
			Expect(res1.Expression.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("owner"))

			ra2 := repositorySt.PermissionStatements[1].(*ast.PermissionStatement)
			Expect(ra2.Name.Literal).Should(Equal("delete"))

			res2 := ra2.ExpressionStatement.(*ast.ExpressionStatement)

			Expect(res2.Expression.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("owner"))
			Expect(res2.Expression.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("parent.create_repository"))
		})

		It("Case 7 - Multiple actions", func() {
			pr := NewParser(`
		entity user {}

		entity organization {
			// relations
			relation admin @user
			relation member @user

			// actions
			action create_repository = (admin or member)
			action manage_team = (admin)
		}

		entity team {
			// relations
			relation leader @user
			relation member @user

			// actions
			permission add_member = (leader or (parent.manage_team))
		}
		`)

			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())

			// USER
			userSt := schema.Statements[0].(*ast.EntityStatement)
			Expect(userSt.Name.Literal).Should(Equal("user"))

			// ORGANIZATION
			organizationSt := schema.Statements[1].(*ast.EntityStatement)
			Expect(organizationSt.Name.Literal).Should(Equal("organization"))

			oa2 := organizationSt.PermissionStatements[1].(*ast.PermissionStatement)
			Expect(oa2.Name.Literal).Should(Equal("manage_team"))

			oes2 := oa2.ExpressionStatement.(*ast.ExpressionStatement)
			Expect(oes2.Expression.(*ast.Identifier).String()).Should(Equal("admin"))

			// TEAM
			teamSt := schema.Statements[2].(*ast.EntityStatement)
			Expect(teamSt.Name.Literal).Should(Equal("team"))

			tperm1 := teamSt.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(tperm1.Name.Literal).Should(Equal("add_member"))

			tes1 := tperm1.ExpressionStatement.(*ast.ExpressionStatement)
			Expect(tes1.Expression.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("leader"))
			Expect(tes1.Expression.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("parent.manage_team"))
		})

		It("Case 8 - Complex nested expressions", func() {
			pr := NewParser(`
	entity user {}

	entity organization {
		// relations
		relation admin @user
		relation member @user

		// actions
		action manage_organization = ((admin and not member) or (member and not admin))
	}

	entity team {
		// relations
		relation leader @user
		relation member @user

		// actions
		permission manage_team = ((leader and parent.manage_organization) or (member and not parent.manage_organization))
	}
	`)

			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())

			// USER
			userSt := schema.Statements[0].(*ast.EntityStatement)
			Expect(userSt.Name.Literal).Should(Equal("user"))

			// ORGANIZATION
			organizationSt := schema.Statements[1].(*ast.EntityStatement)
			Expect(organizationSt.Name.Literal).Should(Equal("organization"))

			oa1 := organizationSt.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(oa1.Name.Literal).Should(Equal("manage_organization"))

			oes1 := oa1.ExpressionStatement.(*ast.ExpressionStatement)
			Expect(oes1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("admin"))
			Expect(oes1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("not member"))
			Expect(oes1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("member"))
			Expect(oes1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("not admin"))

			// TEAM
			teamSt := schema.Statements[2].(*ast.EntityStatement)
			Expect(teamSt.Name.Literal).Should(Equal("team"))

			tperm1 := teamSt.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(tperm1.Name.Literal).Should(Equal("manage_team"))

			tes1 := tperm1.ExpressionStatement.(*ast.ExpressionStatement)
			Expect(tes1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("leader"))
			Expect(tes1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("parent.manage_organization"))
			Expect(tes1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("member"))
			Expect(tes1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("not parent.manage_organization"))
		})

		It("Case 9 - More complex nested expressions", func() {
			pr := NewParser(`
	entity user {}

	entity organization {
		// relations
		relation admin @user
		relation member @user

		// actions
		action manage_organization = (((admin and not member) or member) and (not admin and not member))
	}

	entity project {
		// relations
		relation owner @user
		relation contributor @user

		// actions
		permission manage_project = ((owner and (parent.admin or parent.member)) or (contributor and not parent.manage_organization and (not parent.admin and not parent.member)))
	}
	`)

			schema, err := pr.Parse()
			Expect(err).ShouldNot(HaveOccurred())

			// USER
			userSt := schema.Statements[0].(*ast.EntityStatement)
			Expect(userSt.Name.Literal).Should(Equal("user"))

			// ORGANIZATION
			organizationSt := schema.Statements[1].(*ast.EntityStatement)
			Expect(organizationSt.Name.Literal).Should(Equal("organization"))

			oa1 := organizationSt.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(oa1.Name.Literal).Should(Equal("manage_organization"))

			oes1 := oa1.ExpressionStatement.(*ast.ExpressionStatement)
			Expect(oes1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("admin"))
			Expect(oes1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("not member"))
			Expect(oes1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("member"))
			Expect(oes1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("not admin"))
			Expect(oes1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("not member"))

			// PROJECT
			projectSt := schema.Statements[2].(*ast.EntityStatement)
			Expect(projectSt.Name.Literal).Should(Equal("project"))

			p1 := projectSt.PermissionStatements[0].(*ast.PermissionStatement)
			Expect(p1.Name.Literal).Should(Equal("manage_project"))

			eps1 := p1.ExpressionStatement.(*ast.ExpressionStatement)
			Expect(eps1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("owner"))
			Expect(eps1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("parent.admin"))
			Expect(eps1.Expression.(*ast.InfixExpression).Left.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("parent.member"))
			Expect(eps1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("contributor"))
			Expect(eps1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("not parent.admin"))
			Expect(eps1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.InfixExpression).Left.(*ast.Identifier).String()).Should(Equal("not parent.admin"))
			Expect(eps1.Expression.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.InfixExpression).Right.(*ast.Identifier).String()).Should(Equal("not parent.member"))
		})

		It("Case 10 - Duplicate entity", func() {
			pr := NewParser(`
        entity user {}
        entity user {}
    `)

			_, err := pr.Parse()

			// Ensure an error is returned
			Expect(err).Should(HaveOccurred())

			// Ensure the error message contains the expected string
			Expect(err.Error()).Should(ContainSubstring("3:23:duplication found for user"))
		})

		It("Case 11 - Duplicate Relation", func() {
			pr := NewParser(`
 				entity organization { 
					relation member @user 
					relation member @user 
				} `)

			_, err := pr.Parse()

			// Ensure an error is returned
			Expect(err).Should(HaveOccurred())

			// Ensure the error message contains the expected string
			Expect(err.Error()).Should(ContainSubstring("5:2:duplication found for organization#member"))
		})

		It("Case 12 - Duplicate action", func() {
			pr := NewParser(`
			entity organization {
				relation admin @user
				action delete = admin 
				permission delete = admin 
			}`)

			_, err := pr.Parse()

			// Ensure an error is returned
			Expect(err).Should(HaveOccurred())

			// Ensure the error message contains the expected string
			Expect(err.Error()).Should(ContainSubstring("5:25:duplication found for organization#delete"))
		})
	})
})
