package packer

import (
	"cgl.tideland.biz/asserts"
	"testing"
)

type hashBuilderFactory struct {
	builderMap map[string]Builder
}

func (bf *hashBuilderFactory) CreateBuilder(name string) Builder {
	return bf.builderMap[name]
}

type TestBuilder struct {
	prepareCalled bool
	prepareConfig interface{}
	runCalled bool
	runBuild *Build
	runUi Ui
}

func (tb *TestBuilder) Prepare(config interface{}) {
	tb.prepareCalled = true
	tb.prepareConfig = config
}

func (tb *TestBuilder) Run(b *Build, ui Ui) {
	tb.runCalled = true
	tb.runBuild = b
	tb.runUi = ui
}

func testBuild() *Build {
	return &Build{
		name: "test",
		builder: &TestBuilder{},
		rawConfig: 42,
	}
}

func testBuilder() *TestBuilder {
	return &TestBuilder{}
}

func testBuildFactory(builderMap map[string]Builder) BuilderFactory {
	return &hashBuilderFactory{builderMap}
}

func TestBuild_Prepare(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	build := testBuild()
	builder := build.builder.(*TestBuilder)

	build.Prepare()
	assert.True(builder.prepareCalled, "prepare should be called")
	assert.Equal(builder.prepareConfig, 42, "prepare config should be 42")
}

func TestBuild_Run(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	ui := testUi()

	build := testBuild()
	build.Prepare()
	build.Run(ui)

	builder := build.builder.(*TestBuilder)

	assert.True(builder.runCalled, "run should be called")
	assert.Equal(builder.runBuild, build, "run should be called with build")
	assert.Equal(builder.runUi, ui, "run should be called with ui")
}

func TestBuild_RunBeforePrepare(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	defer func() {
		p := recover()
		assert.NotNil(p, "should panic")
		assert.Equal(p.(string), "Prepare must be called first", "right panic")
	}()

	testBuild().Run(testUi())
}