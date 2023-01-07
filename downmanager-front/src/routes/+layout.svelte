<script lang="ts">
    import TopAppBar , {Row, Section, Title} from "@smui/top-app-bar";
</script>

<div class="main">
    <div class="head">
        <TopAppBar>
            <Row>
                <Section>
                    <Title>Downmanager</Title>
                </Section>
            </Row>
        </TopAppBar>
    </div>
    <div class="content">
        <slot></slot>
    </div>
</div>

<style lang="scss">
    .main {
        margin: -8px;
        display: flex;
        flex-wrap: wrap;
        .head {
 
        }
        .content {
            margin-top: 64px;
            padding: 8px;
        }
    }
</style>